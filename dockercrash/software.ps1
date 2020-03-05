###
### list of files to download
###

## files downloaded from web
$rubyinstaller = 'https://dl.bintray.com/oneclick/rubyinstaller/rubyinstaller-2.3.3-x64.exe'
$rubydevkit = 'https://dl.bintray.com/oneclick/rubyinstaller/DevKit-mingw64-64-4.7.2-20130224-1432-sfx.exe'
$pythonmsi = 'https://www.python.org/ftp/python/2.7.13/python-2.7.13.amd64.msi'
$gomsi = 'https://dl.google.com/go/go1.10.3.windows-amd64.msi'
$wixmsi = 'https://github.com/wixtoolset/wix3/releases/download/wix3111rtm/wix311.exe'
$vcpmsi = 'https://download.microsoft.com/download/7/9/6/796EF2E4-801B-4FC4-AB28-B59FBF6D907B/VCForPython27.msi'
$sevenzip='https://www.7-zip.org/a/7z1801-x64.exe'
$vs_buildtools ='https://aka.ms/vs/15/release/vs_buildtools.exe'
$s3cmd = 'https://s3.amazonaws.com/aws-cli/AWSCLI64.msi'
$git = 'https://github.com/git-for-windows/git/releases/download/v2.19.0.windows.1/Git-2.19.0-64-bit.exe'

## files downloaded from our bucket

## currently these are not used until we can figure out how to get the creds into
## the container
$netfx = 's3://dd-builder-windows-build-unstable/microsoft-windows-netfx3-ondemand-package.cab'
$winbuilds = 's3://dd-builder-windows-build-unstable/winbuilds.zip'
$mingit = 's3://dd-builder-windows-build-unstable/mingit.zip'
$patch = 's3://dd-builder-windows-build-unstable/patch-2.5.9-7-bin.zip'