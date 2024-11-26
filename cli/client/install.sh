#!/usr/bin/env bash

device=$1
mount_point=$2

if [ -z "$device" ]; then
  echo "You must specify a device!" 1>&2
  exit 1
fi

if [ -z "$mount_point" ]; then
  mount_point="/mnt/nf6_install"
fi

sgdisk --zap-all "$device"
sgdisk -n 0:0:+550MiB -t 0:ef00 "$device"
sgdisk -n 0:0:0 -t 0:8300 "$device"

boot_device="/dev/$(lsblk -J "$device" | jq -r ".blockdevices[0].children[0].name")"
root_device="/dev/$(lsblk -J "$device" | jq -r ".blockdevices[0].children[1].name")"

mkfs.fat -F 32 "$boot_device"
mkfs.ext4 -F "$root_device"
mkdir -p "$mount_point"
mount "$root_device" "$mount_point"
mkdir -p "$mount_point/boot"
mount "$boot_device" "$mount_point/boot"

nixos-generate-config --root "$mount_point"
