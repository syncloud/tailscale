local name = 'tailscale';
local headscale = '0.29.7';
local tailscale = '1.104.7';
local go = '1.27.3';
local node = '22-bookworm';
local playwright = 'mcr.microsoft.com/playwright:v1.59.1-jammy';
local python = '3.12-slim-bookworm';
local debian = 'bookworm-slim';
local platform = '26.10.02';
local store_publisher = 'stable-366';
local distros = ['bookworm', 'buster'];
local snapd_hold = "mkdir -p /etc/systemd/system/snapd.service.d && printf '[Service]\\nExecStartPost=/bin/sh -c \"/usr/bin/snap set system refresh.hold=2099-01-01T00:00:00Z\"\\n' > /etc/systemd/system/snapd.service.d/disable-refresh.conf && exec /sbin/init";

local platform_image(distro) = 'syncloud/platform-' + distro + ':' + platform;

local e2e(project) = {
  name: 'e2e ' + project,
  image: playwright,
  commands: ['./ci/ui.sh ' + project + ' bookworm specs/01-devices.spec.ts'],
};

local build(arch, test_ui) = [{
  kind: 'pipeline',
  type: 'docker',
  name: arch,
  platform: {
    os: 'linux',
    arch: arch,
  },
  steps: [
    {
      name: 'headscale',
      image: 'debian:' + debian,
      commands: ['./headscale/build.sh ' + headscale],
    },
    {
      name: 'tailscale',
      image: 'debian:' + debian,
      commands: ['./tailscale/build.sh ' + tailscale],
    },
  ] + [
    {
      name: 'binaries test ' + distro,
      image: platform_image(distro),
      commands: ['./headscale/test.sh', './tailscale/test.sh'],
    }
    for distro in distros
  ] + [
    {
      name: 'web',
      image: 'node:' + node,
      commands: ['./web/build.sh'],
    },
    {
      name: 'cli',
      image: 'golang:' + go,
      commands: ['./cli/build.sh'],
    },
    {
      name: 'backend',
      image: 'golang:' + go,
      commands: ['./backend/build.sh'],
    },
  ] + [
    {
      name: 'cli test ' + distro,
      image: platform_image(distro),
      commands: ['./cli/test.sh'],
    }
    for distro in distros
  ] + [
    {
      name: 'package',
      image: 'debian:' + debian,
      commands: ['./package.sh ' + name + ' $DRONE_BUILD_NUMBER'],
    },
  ] + [
    {
      name: 'test ' + distro,
      image: 'python:' + python,
      commands: ['./test/run.sh ' + distro + ' ' + name + ' $DRONE_BUILD_NUMBER ' + tailscale],
    }
    for distro in distros
  ] + (if test_ui then [e2e('desktop'), e2e('mobile')] else []) + [
    {
      name: 'publish',
      image: 'syncloud/store-publisher:' + store_publisher,
      environment: {
        SYNCLOUD_TOKEN: { from_secret: 'SYNCLOUD_TOKEN' },
      },
      command: ['snap', '-c', '${DRONE_BRANCH}'],
      when: {
        branch: ['stable'],
        event: ['push'],
      },
    },
    {
      name: 'artifact',
      image: 'appleboy/drone-scp:1.6.4',
      settings: {
        host: { from_secret: 'artifact_host' },
        username: 'artifact',
        key: { from_secret: 'artifact_key' },
        timeout: '2m',
        command_timeout: '2m',
        target: '/home/artifact/repo/' + name + '/${DRONE_BUILD_NUMBER}-' + arch,
        source: ['artifact/*'],
        strip_components: 1,
      },
      when: {
        status: ['failure', 'success'],
      },
    },
  ],
  trigger: {
    event: ['push'],
  },
  services: [
    {
      name: name + '.' + distro + '.com',
      image: platform_image(distro),
      privileged: true,
      entrypoint: ['/bin/sh', '-c', snapd_hold],
      volumes: [
        { name: 'dbus', path: '/var/run/dbus' },
        { name: 'dev', path: '/dev' },
      ],
    }
    for distro in distros
  ],
  volumes: [
    { name: 'dbus', host: { path: '/var/run/dbus' } },
    { name: 'dev', host: { path: '/dev' } },
  ],
}];

build('amd64', true) +
build('arm64', false)
