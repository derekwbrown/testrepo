. .\tls12.ps1
. .\software.ps1

Write-Host -ForegroundColor Green Installing Git
$out = 'gitwin.exe'
(New-Object System.Net.WebClient).DownloadFile($git, $out)
Start-Process $out -ArgumentList '/VERYSILENT /DIR="c:\devtools\git"' -Wait

Remove-Item $out
# set path locally so we can initialize git config
$Env:Path="$Env:Path;c:\devtools\git\cmd"
& 'git.exe' config --global user.name "Croissant Builder"
& 'git.exe' config --global user.email "croissant@datadoghq.com"

Write-Host -ForegroundColor Green Done with Git
