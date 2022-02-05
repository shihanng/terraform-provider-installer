resource "installer_script" "this" {
  path           = "/tmp/installer-myapp-test"
  install_script = <<-EOF
  /bin/bash

  touch /tmp/installer-myapp-test
  chmod +x /tmp/installer-myapp-test
  exit 0
  EOF

  uninstall_script = <<-EOF
  /bin/bash

  rm -f /tmp/installer-myapp-test
  exit 0
  EOF
}
