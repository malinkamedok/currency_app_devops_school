import testinfra  # type: ignore
from testinfra.host import Host  # type: ignore

def test_packages(host: Host):
    packages = ['curl', 'openjdk-17-jdk-headless']
    for package in packages:
        assert host.package(package).is_installed