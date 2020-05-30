
# soba: backup hosted git repositories
[![Build Status](https://travis-ci.org/jonhadfield/soba.svg?branch=master)](https://travis-ci.org/jonhadfield/soba) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/8b6afc5274e84d50bd5345580d3a0405)](https://www.codacy.com/app/jonhadfield/soba?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=jonhadfield/soba&amp;utm_campaign=Badge_Grade) [![Codacy Badge](https://api.codacy.com/project/badge/Coverage/8b6afc5274e84d50bd5345580d3a0405)](https://www.codacy.com/app/jonhadfield/soba?utm_source=github.com&utm_medium=referral&utm_content=jonhadfield/soba&utm_campaign=Badge_Coverage) [![Go Report Card](https://goreportcard.com/badge/github.com/jonhadfield/soba)](https://goreportcard.com/report/github.com/jonhadfield/soba)

- [about](#about)
- [configuration](#configuration)
- [run using command line](#run-using-command-line)
- [run using docker](#run-using-docker)
- [run on Synology NAS](#run-on-synology-nas)

## about

soba is tool for backing up private and public git repositories hosted on the most popular hosting providers. It generates a [git bundle](https://git-scm.com/book/en/v2/Git-Tools-Bundling) that stores a backup of each repository as a single file. 

An unchanged git repository will create an identical bundle file so bundles will only be stored if a change has been made and will not produce duplicates.

#### Supported OSes

Tested on Windows 10, MacOS, and Linux (amd64).   
Not tested, but should also work on builds for: Linux (386, arm386 and arm64), FreeBSD, NetBSD, and OpenBSD.

#### Supported Providers

- BitBucket
- GitHub
- GitLab

## configuration

soba can be run from the command line or as docker container. In both cases the only configuration required is an environment variable with the directory in which to create backups, and additional to define credentials for each the providers. 

On Windows 10: 
- search for 'environment variables' and choose 'Edit environment variables for your account'
- choose 'New...' under the top pane and enter the name/key and value for each of the settings

On Linux and MacOS you would set these using:

```bash
$ export GIT_BACKUP_DIR="/repo-backups/"
```

To set provider tokens see [below](#setting-provider-tokens).

### run using command line

Download the latest release [here](https://github.com/jonhadfield/soba/releases) and then install:
```
$ install <soba binary> /usr/local/bin/soba
```

After setting GIT_BACKUP_DIR, set your provider token(s) as detailed [here](#setting-provider-tokens).

and then run:

```bash
$ soba
```

### run using docker

Using docker enables you to run soba without anything else installed.

Docker requires you pass environment variables to the container using the '-e' option and that you mount your preferred backup directory. For example:

```bash
$ docker run --rm -t \
             -v <your backup dir>:/backup \
             -e GIT_BACKUP_DIR='/backup' \
             -e GITHUB_TOKEN='MYGITHUBTOKEN' \
             -e GITLAB_TOKEN='MYGITLABTOKEN' \
             quay.io/jeffplourde/soba
```

To hide credentials, you can instead use exported environment variables and specify using this syntax:

```bash
$ docker run --rm -t \
             -v <your backup dir>:/backup \
             -e GIT_BACKUP_DIR='/backup' \
             -e GITHUB_TOKEN=$GITHUB_TOKEN \
             -e GITLAB_TOKEN=$GITLAB_TOKEN \
             quay.io/jeffplourde/soba
```


## scheduling backups

Backups can be scheduled to run by setting an additional environment variable: GIT_BACKUP_INTERVAL. The value is the number of hours between backups. For example, this will run the backup daily:

```bash
$ export GIT_BACKUP_INTERVAL=24
```

if using docker then add:

```bash
-e GIT_BACKUP_INTERVAL=24
```

_Note: the interval is added to the start of the last backup and not the time it finished. Therefore, ensure the interval is greater than the duration of a backup._

## setting provider tokens

On Linux and MacOS you can set environment variables manually before each time you run soba:

```bash
$ export NAME='VALUE'
```
    
or by defining in a startup file for your shell so they are automatically set and available when you need them. For example, if using the bash shell and running soba as your user, add the relevant export statements to the following file: 

```
/home/<your-user-id>/.bashrc
```

and run:

```bash
$ source /home/<your-user-id>/.bashrc
```

| Provider | Environment Variable(s) | Generating token |
|:---------|:---------------|:----------------------------------------------------------------------------------------------------------------------------------|
| BitBucket| BITBUCKET_USER<br/>BITBUCKET_KEY<br/>BITBUCKET_SECRET | <a href="https://confluence.atlassian.com/bitbucket/oauth-on-bitbucket-cloud-238027431.html" target="_blank">instructions</a>
| GitHub   | GITHUB_TOKEN   | <a href="https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/" target="_blank">instructions</a>
| GitLab   | GITLAB_TOKEN   | <a href="https://gitlab.com/profile/personal_access_tokens" target="_blank">instructions</a>

### run on Synology NAS
_Tested on DS916+_


1. Create a directory on your NAS for backing up Git repositories to
2. Install Docker from the Synology Package Center
3. Open Docker and select 'Image'
4. Select 'Add' from the top menu and choose 'Add From URL'
5. In 'Repository URL' enter 'jeffplourde/soba', leave other options as default and click 'Add'
6. When it asks to 'Choose Tag' accept the default 'latest' by pressing 'Select'
7. Select image 'jeffplourde/soba:latest' from the list and click 'Launch' from the top menu
8. Set 'Container Name' to 'soba' and select 'Advanced Settings'
9. Check 'Enable auto-restart'
10. Under 'Volume' select 'Add folder' and choose the directory created in step 1. Set the 'Mount Path' to '/backup'
11. Under 'Network' check 'Use the same network as Docker Host'
12. Under 'Environment' click '+' to add each of the following:
  - **variable** GIT_BACKUP_DIR **Value** /backup
  - **variable** GIT_BACKUP_INTERVAL **Value** [hours between backups]
#### Provider Specific
  - **variable** BITBUCKET_USER **Value** [BitBucket User]   (if using BitBucket)
  - **variable** BITBUCKET_KEY **Value** [BitBucket Key]   (if using BitBucket)
  - **variable** BITBUCKET_SECRET **Value** [BitBucket Secret]   (if using BitBucket)
  - **variable** GITHUB_TOKEN **Value** [GitHub token]   (if using GitHub)
  - **variable** GITLAB_TOKEN **Value** [GitLab token]   (if using GitLab)

13. Click 'Apply'
14. Leave settings as default and select 'Next'
15. Check 'Run this container after the wizard is finished' and click 'Apply'

The container should launch in a few seconds. You can view progress by choosing 'Container' in the left-hand menu, select 'soba', choose 'details' and then click on 'Log'
  

