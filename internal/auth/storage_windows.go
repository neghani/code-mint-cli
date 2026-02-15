//go:build windows

package auth

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type windowsStore struct {
	target   string
	fallback *fileStore
}

func newPlatformStore(profile string) (TokenStore, error) {
	fs, err := newFileStore(profile)
	if err != nil {
		return nil, err
	}
	return &windowsStore{target: "codemint-" + profile, fallback: fs}, nil
}

func (w *windowsStore) Set(ctx context.Context, token string) error {
	out, err := exec.Command("cmdkey", "/generic:"+w.target, "/user:codemint", "/pass:"+token).CombinedOutput()
	if err == nil {
		return nil
	}
	if ferr := w.fallback.Set(ctx, token); ferr == nil {
		return nil
	}
	return fmt.Errorf("write windows credential manager token: %s", strings.TrimSpace(string(out)))
}

func (w *windowsStore) Get(ctx context.Context) (string, error) {
	token, err := w.readCredential()
	if err == nil && token != "" {
		return token, nil
	}
	return w.fallback.Get(ctx)
}

func (w *windowsStore) Delete(ctx context.Context) error {
	_, _ = exec.Command("cmdkey", "/delete:"+w.target).CombinedOutput()
	_ = w.fallback.Delete(ctx)
	return nil
}

func (w *windowsStore) readCredential() (string, error) {
	script := `
Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;
using System.Text;

public class NativeCred {
  [StructLayout(LayoutKind.Sequential, CharSet=CharSet.Unicode)]
  public struct CREDENTIAL {
    public UInt32 Flags;
    public UInt32 Type;
    public string TargetName;
    public string Comment;
    public System.Runtime.InteropServices.ComTypes.FILETIME LastWritten;
    public UInt32 CredentialBlobSize;
    public IntPtr CredentialBlob;
    public UInt32 Persist;
    public UInt32 AttributeCount;
    public IntPtr Attributes;
    public string TargetAlias;
    public string UserName;
  }

  [DllImport("advapi32.dll", EntryPoint="CredReadW", CharSet=CharSet.Unicode, SetLastError=true)]
  public static extern bool CredRead(string target, int type, int reservedFlag, out IntPtr credentialPtr);

  [DllImport("advapi32.dll", EntryPoint="CredFree", SetLastError=true)]
  public static extern void CredFree([In] IntPtr cred);
}
"@

$target = %s
$ptr = [IntPtr]::Zero
if (-not [NativeCred]::CredRead($target, 1, 0, [ref]$ptr)) { exit 1 }
try {
  $cred = [Runtime.InteropServices.Marshal]::PtrToStructure($ptr, [Type][NativeCred+CREDENTIAL])
  if ($cred.CredentialBlobSize -le 0) { exit 1 }
  $bytes = New-Object byte[] $cred.CredentialBlobSize
  [Runtime.InteropServices.Marshal]::Copy($cred.CredentialBlob, $bytes, 0, $cred.CredentialBlobSize)
  [Text.Encoding]::Unicode.GetString($bytes)
} finally {
  [NativeCred]::CredFree($ptr)
}
`
	command := fmt.Sprintf(script, psLiteral(w.target))
	out, err := exec.Command("powershell", "-NoProfile", "-Command", command).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("read windows credential manager token: %s", strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func psLiteral(v string) string {
	return "'" + strings.ReplaceAll(v, "'", "''") + "'"
}
