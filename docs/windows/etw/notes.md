# Notes

## Microsoft-Windows-Kernel-Process

The following keywords are available:

- WINEVENT_KEYWORD_PROCESS
- WINEVENT_KEYWORD_THREAD
- WINEVENT_KEYWORD_IMAGE
- WINEVENT_KEYWORD_CPU_PRIORITY
- WINEVENT_KEYWORD_OTHER_PRIORITY
- WINEVENT_KEYWORD_PROCESS_FREEZE
- WINEVENT_KEYWORD_JOB
- WINEVENT_KEYWORD_ENABLE_PROCESS_TRACING_CALLBACKS
- WINEVENT_KEYWORD_JOB_IO
- WINEVENT_KEYWORD_WORK_ON_BEHALF
- WINEVENT_KEYWORD_JOB_SILO

See:

- [repnz/etw-providers-docs - Win10-17134 ETW providers](https://github.com/repnz/etw-providers-docs/blob/master/Manifests-Win10-17134/Microsoft-Windows-Kernel-Process.xml)

## PID -> PPID

The following events contain both a PID and PPID:

```xml
...
<template tid="ProcessStartArgs">
  <data name="ProcessID" inType="win:UInt32"/>
  <data name="CreateTime" inType="win:FILETIME"/>
  <data name="ParentProcessID" inType="win:UInt32"/>
  <data name="SessionID" inType="win:UInt32"/>
  <data name="ImageName" inType="win:UnicodeString"/>
</template>
...
<template tid="ProcessStartArgs_V1">
  <data name="ProcessID" inType="win:UInt32"/>
  <data name="CreateTime" inType="win:FILETIME"/>
  <data name="ParentProcessID" inType="win:UInt32"/>
  <data name="SessionID" inType="win:UInt32"/>
  <data name="Flags" inType="win:UInt32"/>
  <data name="ImageName" inType="win:UnicodeString"/>
</template>
...
<template tid="ProcessStartArgs_V2">
  <data name="ProcessID" inType="win:UInt32"/>
  <data name="CreateTime" inType="win:FILETIME"/>
  <data name="ParentProcessID" inType="win:UInt32"/>
  <data name="SessionID" inType="win:UInt32"/>
  <data name="Flags" inType="win:UInt32"/>
  <data name="ImageName" inType="win:UnicodeString"/>
  <data name="ImageChecksum" inType="win:UInt32"/>
  <data name="TimeDateStamp" inType="win:UInt32"/>
  <data name="PackageFullName" inType="win:UnicodeString"/>
  <data name="PackageRelativeAppId" inType="win:UnicodeString"/>
</template>
...
```
