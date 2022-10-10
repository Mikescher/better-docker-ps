# ./dops - better `docker ps` 
A replacement for the default docker-ps that tries really hard to fit into the width of your terminal°

![](readme.d/default.png)

## Rationale

By default, my `docker ps` output is really wide and every line wraps around into three.
This (obviously) breaks the tabular display and makes everything chaotic.  
It does not look like we will get a better output in the foreseeable future (see [moby#7477](https://github.com/moby/moby/issues/7477)), so I decided to make my own drop-in replacement.  

## Features

 - All normal commandline flags/options from docker-ps work *(almost)* the same.
 - Write multi-value data (like multiple port mappings, multiple networks, etc.) into multiple lines instead of concatenating them.
 - Add color to the STATE and STATUS column (green / yellow / red).
 - Automatically remove columns in the output until it fits in the current terminal width.


More Changes from default docker-ps:
 - Show (by default) the container-cmd without arguments.
 - Show the ImageName (by default) without the registry prefix, and split ImageName and ImageTag into two columns.
 - Added the columns IP and NETWORK to the default column set (if they fit)
 - Added support for a few new columns (via --format):  
   `{{.ImageName}`, `{{.ImageTag}`, `{{.Tag}`, `{{.ImageRegistry}`, `{{.Registry}`, `{{.ShortCommand}`, `{{.LabelKeys}`, `{{.IP}`                         
 - Added options to control the color-output, the used socket, the time-zone and time-format, etc (see `./dops --help`) 

## Getting started

 - Download the latest binary from the [releases page](https://github.com/Mikescher/better-docker-ps/releases)
 - but it into yout PATH (eg /usr/local/bin)
 - (optional) alias teh docker ps command (see [section below](#usage-as-drop-in-replacement))

## Screenshots

![](readme.d/fullsize.png)  
All (default) columns visible

&nbsp;

![](readme.d/default.png)  
Output on a medium sized terminal

&nbsp;

![](readme.d/small.png)  
Output on a small terminal

&nbsp;

## Usage as drop-in replacement

You can fully replace docker ps by creating a shell function in your `.bashrc` / `.zshrc`...

~~~sh
docker() {
  case $1 in
    ps)
      shift
      command dops "$@"
      ;;
    *)
      command docker "$@";;
  esac
}
~~~

This will alias every call to `docker ps ...` with `dops ...` (be sure to have the dops binary in your PATH).

If you are using the fish-shell you have to create a (similar) function:

~~~fish
function docker
    if test -n "$argv[1]"
        switch $argv[1]
            case ps
                dops $argv[2..-1]
            case '*'
                command docker $argv[1..-1]
        end
    end
end
~~~

## Manual

Output of `./dops --help`:

~~~~~~
better-docker-ps

Usage:
  dops [OPTIONS]                     List docker container

Options (default):
  -h, --help                         Show this screen.
  --version                          Show version.
  --all , -a                         Show all containers (default shows just running)
  --filter <ftr>, -f <ftr>           Filter output based on conditions provided
  --format <fmt>                     Pretty-print containers using a Go template
  --last , -n                        Show n last created containers (includes all states)
  --latest , -l                      Show the latest created container (includes all states)
  --no-trunc                         Don't truncate output
  --quiet , -q                       Only display container IDs
  --size , -s                        Display total file sizes

Options (extra | do not exist in `docker ps`):
  --silent                           Do not print any output
  --timezone                         Specify the timezone for date outputs
  --color <true|false>               Enable/Disable terminal color output
  --no-color                         Disable terminal color output
  --socket <filepath>                Specify the docker socket location (Default: /var/run/docker.sock)
  --timeformat <go-time-fmt>         Specify the datetime output format (golang syntax)
  --no-header                        Do not print the table header
  --simple-header                    Do not print the lines under the header
  --format <fmt>                     You can specify multiple formats and the first one that fits your terminal widt will be used

Available --format keys (default):
  {{.ID}                             Container ID
  {{.Image}                          Image ID
  {{.Command}                        Quoted command
  {{.CreatedAt}                      Time when the container was created.
  {{.RunningFor}                     Elapsed time since the container was started.
  {{.Ports}                          Exposed ports.
  {{.State}                          Container status
  {{.Status}                         Container status with details
  {{.Size}                           Container disk size.
  {{.Names}                          Container names.
  {{.Labels}                         All labels assigned to the container.
  {{.Label}                          [!] Unsupported
  {{.Mounts}                         Names of the volumes mounted in this container.
  {{.Networks}                       Names of the networks attached to this container.

Available --format keys (extra | do not exist in `docker ps`):
  {{.ImageName}                      Image ID (without tag and registry)
  {{.ImageTag}, {{.Tag}              Image Tag
  {{.ImageRegistry}, {{.Registry}    Image Registry
  {{.ShortCommand}                   Command without arguments
  {{.LabelKeys}                      All labels assigned to the container (keys only)
  {{.IP}                             Internal IP Address
~~~~~~
