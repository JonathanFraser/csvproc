FROM pachyderm/job-shim
MAINTAINER jonathan.fraser@generalfusion.com
ADD convertcsv/convertcsv /convertcsv
ADD folder_process.sh folder_process.sh
