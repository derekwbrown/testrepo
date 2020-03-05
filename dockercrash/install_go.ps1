. .\tls12.ps1
. .\software.ps1

Write-Host -ForegroundColor Green Installing Go
$out = 'go.msi'
(New-Object System.Net.WebClient).DownloadFile($gomsi, $out)
Invoke-WebRequest -OutFile go.msi -Uri https://dl.google.com/go/go1.10.3.windows-amd64.msi 
Start-Process msiexec -ArgumentList '/q /i go.msi' -Wait
Remove-Item $out

md c:\dev\go
setx GOPATH c:\dev\go
setx PATH "$Env:Path;c:\dev\go\bin"
Write-Host -ForegroundColor Green Done with Go
