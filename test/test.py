import json
import os
import subprocess
import time
from os.path import dirname, join
from subprocess import check_output

import pytest
import requests
from syncloudlib.integration.hosts import add_host_alias
from syncloudlib.integration.installer import local_install, wait_for_installer

import tunnel

DIR = dirname(__file__)
TMP_DIR = '/tmp/syncloud'
CLIENT_DIR = '/tmp/tailscale-client'
CLIENT_STATE = '/tmp/tailscale-client-state'
CLIENT_SOCKET = join(CLIENT_STATE, 'tailscaled.sock')
SOCKS_PORT = 1055
HEADSCALE = '/snap/tailscale/current/bin/headscale.sh'
SERVICES = ['headscale', 'tailscaled', 'backend']


@pytest.fixture(scope="session")
def module_setup(request, device, platform_data_dir, app_dir, artifact_dir, snap_data_dir):
    def module_teardown():
        device.run_ssh('mkdir -p {0}'.format(TMP_DIR), throw=False)
        device.run_ssh('journalctl > {0}/journalctl.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ps auxfw > {0}/ps.log'.format(TMP_DIR), throw=False)
        device.run_ssh('{0} nodes list > {1}/nodes.log'.format(HEADSCALE, TMP_DIR), throw=False)
        device.run_ssh('cp {0}/extra-records.json {1}/'.format(snap_data_dir, TMP_DIR), throw=False)
        app_log_dir = join(artifact_dir, 'log')
        os.makedirs(app_log_dir, exist_ok=True)
        device.scp_from_device('{0}/*'.format(TMP_DIR), app_log_dir)
        device.scp_from_device('{0}/config/*.yaml'.format(snap_data_dir), artifact_dir, throw=False)
        if os.path.exists(join(CLIENT_STATE, 'tailscaled.log')):
            check_output('cp {0}/tailscaled.log {1}/client-tailscaled.log'.format(CLIENT_STATE, artifact_dir), shell=True)
        check_output('chmod -R a+r {0}'.format(artifact_dir), shell=True)

    request.addfinalizer(module_teardown)


def settle(device):
    device.run_ssh('snap wait system seed.loaded', retries=100, throw=False)
    device.run_ssh('snap set system refresh.hold=2099-01-01T00:00:00Z', retries=20, throw=False)
    device.run_ssh('snap abort --last=auto-refresh', throw=False)
    device.run_ssh('snap watch --last=auto-refresh', throw=False)


def client(*args, timeout=90):
    return subprocess.run([join(CLIENT_DIR, 'tailscale'), '--socket=' + CLIENT_SOCKET] + list(args),
                          capture_output=True, text=True, timeout=timeout)


def retry(check, attempts=40, delay=5):
    error = None
    for _ in range(attempts):
        try:
            return check()
        except AssertionError as e:
            error = e
            time.sleep(delay)
    raise error


def test_start(module_setup, device, app, domain, device_host):
    add_host_alias(app, device_host, domain)
    device.run_ssh('date', retries=100, throw=True)
    device.run_ssh('mkdir -p {0}'.format(TMP_DIR))
    settle(device)


def test_activate_device(device):
    response = device.activate_custom()
    assert response.status_code == 200, response.text


def test_install(app_archive_path, device_host, device_session, device_password, domain):
    local_install(device_host, device_password, app_archive_path)
    wait_for_installer(device_session, domain)


def test_services_active(device, app):
    def check():
        for service in SERVICES:
            state = device.run_ssh('systemctl is-active snap.{0}.{1}.service'.format(app, service), throw=False)
            assert state.strip() == 'active', '{0}: {1}'.format(service, state)
    retry(check)


def test_health(app_domain):
    def check():
        response = requests.get('https://{0}/health'.format(app_domain), verify=False)
        assert response.status_code == 200, response.text
        assert response.json()['status'] == 'pass', response.text
    retry(check)


def test_info_page_shows_server(app_domain):
    response = requests.get('https://{0}/'.format(app_domain), verify=False)
    assert response.status_code == 200, response.text
    assert 'https://{0}'.format(app_domain) in response.text


def test_derp_reachable_through_platform(app_domain):
    response = requests.get('https://{0}/derp/probe'.format(app_domain), verify=False)
    assert response.status_code == 200, response.text


def test_device_node_online(device, domain):
    hostname = domain.split('.')[0]

    def check():
        nodes = json.loads(device.run_ssh('{0} nodes list -o json'.format(HEADSCALE)))
        node = [n for n in nodes if n['name'] == hostname]
        assert node, nodes
        assert node[0].get('online'), node
        return node[0]
    node = retry(check)
    assert node['ip_addresses'], node


def test_records_point_device_at_node(device, domain, snap_data_dir):
    def check():
        records = json.loads(device.run_ssh('cat {0}/extra-records.json'.format(snap_data_dir)))
        names = [r['name'] for r in records]
        assert domain in names, records
        assert 'auth.{0}'.format(domain) in names, records
        assert 'tailscale.{0}'.format(domain) not in names, records
    retry(check)


def test_client_daemon(device, artifact_dir):
    os.makedirs(CLIENT_STATE, exist_ok=True)
    ca = join(CLIENT_STATE, 'syncloud.ca.crt')
    device.scp_from_device('/var/snap/platform/current/syncloud.ca.crt', ca)
    env = dict(os.environ, SSL_CERT_FILE=ca)
    log = open(join(CLIENT_STATE, 'tailscaled.log'), 'w')
    subprocess.Popen([join(CLIENT_DIR, 'tailscaled'), '--tun=userspace-networking',
                      '--statedir=' + CLIENT_STATE, '--socket=' + CLIENT_SOCKET,
                      '--socks5-server=127.0.0.1:{0}'.format(SOCKS_PORT)],
                     env=env, stdout=log, stderr=subprocess.STDOUT)

    def check():
        result = client('status', '--json')
        assert result.stdout.strip().startswith('{'), result.stderr
    retry(check, attempts=20, delay=2)


def test_interactive_login_goes_to_authelia(app_domain, domain):
    result = client('up', '--login-server=https://{0}'.format(app_domain), '--hostname=interactive', '--timeout=20s')
    output = result.stdout + result.stderr
    urls = [w for w in output.split() if w.startswith('https://{0}/register/'.format(app_domain))]
    assert urls, output
    response = requests.get(urls[0], verify=False, allow_redirects=False)
    assert response.status_code in (302, 303, 307), response.text
    assert 'auth.{0}'.format(domain) in response.headers['Location'], response.headers


def test_client_joins_with_key(device, app_domain):
    users = json.loads(device.run_ssh('{0} users list -o json'.format(HEADSCALE)))
    user = [u for u in users if u['name'] == 'ci']
    if not user:
        user = [json.loads(device.run_ssh('{0} users create ci -o json'.format(HEADSCALE)))]
    key = json.loads(device.run_ssh('{0} preauthkeys create --user {1} --expiration 1h -o json'.format(HEADSCALE, user[0]['id'])))['key']
    result = client('up', '--reset', '--login-server=https://{0}'.format(app_domain), '--authkey=' + key, '--hostname=ci', '--timeout=60s')
    assert result.returncode == 0, result.stdout + result.stderr


def device_ip(device, domain):
    hostname = domain.split('.')[0]
    nodes = json.loads(device.run_ssh('{0} nodes list -o json'.format(HEADSCALE)))
    node = [n for n in nodes if n['name'] == hostname][0]
    return [ip for ip in node['ip_addresses'] if '.' in ip][0]


def test_client_resolves_device_name(device, domain):
    ip = device_ip(device, domain)

    def check():
        result = client('dns', 'query', domain, 'A')
        assert ip in result.stdout, result.stdout + result.stderr
    retry(check, attempts=20)


def test_client_pings_device_directly(domain):
    hostname = domain.split('.')[0]

    def check():
        result = client('ping', '-c', '5', hostname)
        output = result.stdout + result.stderr
        assert 'pong from {0}'.format(hostname) in output, output
        assert 'via DERP' not in output.splitlines()[-1], output
    retry(check, attempts=10)


def test_client_reaches_device_through_tunnel(device, domain, app_domain):
    ip = device_ip(device, domain)

    def check():
        response = tunnel.https_get(SOCKS_PORT, ip, 443, app_domain, '/health')
        assert response.startswith('HTTP/1.1 200'), response
        assert '"pass"' in response, response
    retry(check, attempts=10)
