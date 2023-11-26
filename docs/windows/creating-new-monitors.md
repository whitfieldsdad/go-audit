# Creating new monitors

When creating monitors which source events through Event Tracing for Windows (ETW):

1. Determine which provider to use (e.g. Microsoft-Windows-Kernel-Process)
2. [Lookup the provider GUID](#listing-available-providers) (e.g. `{22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716}` for Microsoft-Windows-Kernel-Process)
3. [Calculate the provider bitmask](#calculate-the-provider-bitmask) using the keywords tied to a class of events provided by a given provider
4. Write a new monitor using the provider GUID and bitmask (see the `ProcessMonitor` struct as an example)

## Listing available providers

```powershell
logman query providers
```

```text
Provider                                 GUID
-------------------------------------------------------------------------------
Microsoft-Windows-COMRuntime             {BF406804-6AFA-46E7-8A48-6C357E1D6D61}
Microsoft-Windows-Kernel-File            {EDD08927-9CC4-4E65-B970-C2560FB5C289}
Microsoft-Windows-Kernel-Network         {7DD42A49-5329-4832-8DFD-43D979153A88}
Microsoft-Windows-Kernel-Process         {22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716}
Microsoft-Windows-Kernel-Registry        {70EB4F03-C1DE-4F73-A051-33D13D5413BD}
Microsoft-Windows-LDAP-Client            {099614A5-5DD7-4788-8BC9-E29F43DB28FC}
Microsoft-Windows-RPC                    {6AD52B32-D609-4BE9-AE07-CE8DAE937E39}
Microsoft-Windows-Services-Svchost       {06184C97-5201-480E-92AF-3A3626C5B140}
Microsoft-Windows-Winlogon               {DBE9B383-7CF3-4331-91CC-A3CB16A3B538}
...
```

> A list of ETW providers extracted from a Windows 11 system is available [here](etw-providers.txt).

## Listing event keywords

```powershell
logman query providers "{22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716}"
```

```text
Provider                                 GUID
-------------------------------------------------------------------------------
Microsoft-Windows-Kernel-Process         {22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716}

Value               Keyword              Description
-------------------------------------------------------------------------------
0x0000000000000010  WINEVENT_KEYWORD_PROCESS
0x0000000000000020  WINEVENT_KEYWORD_THREAD
0x0000000000000040  WINEVENT_KEYWORD_IMAGE
0x0000000000000080  WINEVENT_KEYWORD_CPU_PRIORITY
0x0000000000000100  WINEVENT_KEYWORD_OTHER_PRIORITY
0x0000000000000200  WINEVENT_KEYWORD_PROCESS_FREEZE
0x0000000000000400  WINEVENT_KEYWORD_JOB
0x0000000000000800  WINEVENT_KEYWORD_ENABLE_PROCESS_TRACING_CALLBACKS
0x0000000000001000  WINEVENT_KEYWORD_JOB_IO
0x0000000000002000  WINEVENT_KEYWORD_WORK_ON_BEHALF
0x0000000000004000  WINEVENT_KEYWORD_JOB_SILO
0x8000000000000000  Microsoft-Windows-Kernel-Process/Analytic

Value               Level                Description
-------------------------------------------------------------------------------
0x04                win:Informational    Information

PID                 Image
-------------------------------------------------------------------------------
0x00000000
```

## Calculate the provider bitmask

To calculate the provider bitmask, simply add the value of the keywords you'd like to enable.

For example, to trace both process and image events from the `Microsoft-Windows-Kernel-Process` provider, given:

```text
Value               Keyword              Description
-------------------------------------------------------------------------------
0x0000000000000010  WINEVENT_KEYWORD_PROCESS
0x0000000000000040  WINEVENT_KEYWORD_IMAGE
```

Use the following bitmask:

```text
0x10 + 0x40 = 0x50
```

> When combining hex values with leading zeros, you can simply ignore the leading zeros (i.e. `0x0000000000000010` is the same as `0x10`)

> You can calculate bitmasks using Python: `hex(0x10 + 0x40)` -> `'0x50'`

> [You can calculate bitmasks using the Windows calculator in programmer mode](adding-hex-values-in-windows-calculator.png).

## Additional resources

- https://www.ired.team/miscellaneous-reversing-forensics/windows-kernel-internals/etw-event-tracing-for-windows-101
