#!/bin/sh
repo_name="fake_data_test"
pachctl create-repo $repo_name
for commit in $(seq 1 6); do
	commitid=$(pachctl start-commit $repo_name)
	for f in $(seq 1 25); do
		gencsv | pachctl put-file $repo_name $commitid $f
	done
	pachctl finish-commit $repo_name $commitid
done
