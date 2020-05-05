# ova-forge packer post-processor plugin
A [plugin for packer](https://packer.io/docs/post-processors/index.html)
that allows users to convert VirtualBox ISO builds into OVAs for both
VirtualBox and VMWare - without rebuilding the virtual machine!

## Use cases
Certain VMWare tools (such as vSphere) are very picky about OVF configurations.
As of vSphere 6.5.x, it is not possible to create a VirtualBox OVF and deploy
it without modification. In addition, introducing VMWare into build tooling is
fraught with gotchas and headaches in the forms of licensing, unfixed bugs,
and lack of support.

This plugin effectively allows developers to build one VirtualBox VM and OVF,
and then create VirtualBox and VMWare OVAs. This means that we never need to
touch a VMWare tool to build a VMWare compatible appliance - greatly speeding
up and simplifying the build process.

## How does it work?
After installing the plugin, add the plugin to your packer config's
post-processors section:
```
(...)
    }
  ],
  "post-processors": ["ova-forge"],
  "variables": {
(...)
```

This will convert your VirtualBox OVF into VirtualBox and VMWare OVAs. This
will be indicated in the build output:
```log
==> virtualbox-iso: Exporting virtual machine...
    virtualbox-iso: Executing: export centos7 --output output-centos7-virtualbox-iso/centos7.ovf
==> virtualbox-iso: Deregistering and deleting VM...
==> virtualbox-iso: Running post-processor: ova-forge
    virtualbox-iso (ova-forge): VMWareifying 'output-centos7-virtualbox-iso/centos7.ovf'...
    virtualbox-iso (ova-forge): Finished VMWareifying .ovf at 'output-centos7-virtualbox-iso/centos7-vmware.ovf'
    virtualbox-iso (ova-forge): Creating .ova for 'output-centos7-virtualbox-iso/centos7.ovf' with files 'output-centos7-virtualbox-iso/centos7-disk001.vmdk'...
    virtualbox-iso (ova-forge): Finished creating .ova at 'output-centos7-virtualbox-iso/centos7.ova'
    virtualbox-iso (ova-forge): Creating .ova for 'output-centos7-virtualbox-iso/centos7-vmware.ovf' with files 'output-centos7-virtualbox-iso/centos7-disk001.vmdk'...
    virtualbox-iso (ova-forge): Finished creating .ova at 'output-centos7-virtualbox-iso/centos7-vmware.ova'
Build 'virtualbox-iso' finished.

==> Builds finished. The artifacts of successful builds are:
--> virtualbox-iso: VM files in directory: output-centos7-virtualbox-iso
--> virtualbox-iso: OVA related files: output-centos7-virtualbox-iso/centos7-vmware.ovf, output-centos7-virtualbox-iso/centos7.ova, output-centos7-virtualbox-iso/centos7-vmware.ova
```

## How do I install it?
As of Packer version 1.3.2, you need to do the following:

1. Download the compiled binary for your system from the releases / tags page
2. Create the following directory path in your home directory:
`.packer.d/plugins/`
3. Move the plugin into the directory and make sure it is named
`packer-post-processor-ova-forge`
4. Make sure it is set as executable (on *nix systems)
