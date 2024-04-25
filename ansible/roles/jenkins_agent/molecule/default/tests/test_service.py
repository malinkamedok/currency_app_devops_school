import testinfra  # type: ignore
from testinfra.host import Host  # type: ignore

def test_directories_exist(host):
    directories = ['/usr/local/jenkins-service', '/opt/jenkins']
    for directory in directories:
        assert host.file(directory).exists
        assert host.file(directory).is_directory
        assert host.file(directory).mode == 0o644

def test_swarm_client_downloaded(host):
    assert host.file('/opt/jenkins/swarm-client.jar').exists
    assert host.file('/opt/jenkins/swarm-client.jar').mode == 0o755
