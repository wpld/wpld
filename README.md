# WordPress Local Docker (wpld)

A powerful CLI tool for managing WordPress development environments using Docker. wpld simplifies the process of creating, managing, and working with containerized WordPress installations for local development.

## Features

- 🚀 **Quick Setup**: Create new WordPress projects with a single command
- 🐳 **Docker-based**: Fully containerized development environment
- 🔧 **WP-CLI Integration**: Execute WordPress CLI commands directly
- 🌐 **Network Management**: Automated Docker network configuration
- 📊 **PHPMyAdmin**: Built-in database management interface
- 🔄 **Environment Control**: Easy start, stop, restart operations
- 📝 **Logging**: Access container logs for debugging
- ⚙️ **Configurable**: Flexible configuration options

## Requirements

- [Docker](https://www.docker.com/get-started) installed and running
- [Go](https://golang.org/doc/install) 1.21 or later (for building from source)

## Installation

### Option 1: Install from Source

```bash
# Clone the repository
git clone https://github.com/your-username/wpld.git
cd wpld

# Install the binary
make install
```

### Option 2: Build Binary

```bash
# Build the binary to ./bin/wpld
make build

# Add to your PATH or move to a directory in your PATH
sudo mv bin/wpld /usr/local/bin/
```

## Quick Start

### 1. Create a New WordPress Project

```bash
wpld new
```

This command will:
- Prompt you for project configuration
- Create the necessary Docker containers
- Set up WordPress with proper configuration
- Configure networking and services
- Provide access URLs

### 2. Start an Existing Project

```bash
# Navigate to your project directory
cd my-wordpress-project

# Start the environment
wpld start
```

### 3. Execute WordPress CLI Commands

```bash
# Install a plugin
wpld wp plugin install woocommerce --activate

# Create a user
wpld wp user create admin admin@example.com --role=administrator

# Update WordPress core
wpld wp core update
```

## Commands

### Core Commands

#### `wpld new`
Creates a new WordPress project with interactive setup.

```bash
wpld new
```

#### `wpld start` (alias: `up`)
Starts project services and containers.

```bash
wpld start [flags]

Flags:
  -p, --pull                 Force pulling images before starting containers
  -P, --persist-containers   Do not auto-remove containers on stop
```

#### `wpld stop`
Stops all project containers and services.

```bash
wpld stop
```

#### `wpld restart`
Restarts the project environment (stop + start).

```bash
wpld restart
```

### WordPress Management

#### `wpld wp COMMAND [ARG...]`
Executes WP-CLI commands inside the WordPress container.

```bash
# Examples:
wpld wp core version
wpld wp plugin list
wpld wp theme activate twentytwentyfour
wpld wp db export backup.sql
```

### Utility Commands

#### `wpld exec`
Execute commands inside project containers.

```bash
wpld exec [container] [command]
```

#### `wpld logs`
View container logs for debugging.

```bash
wpld logs [container]
```

#### `wpld run`
Run one-time commands in new containers.

```bash
wpld run [image] [command]
```

#### `wpld config`
Display current project configuration.

```bash
wpld config
```

### Global Options

- `-v, --verbose`: Enable verbose output (can be used multiple times for increased verbosity)
- `--version`: Display version information

## Project Structure

When you create a new project, wpld generates the following structure:

```
my-wordpress-project/
├── .wpld/                 # wpld configuration
│   ├── config.yml        # Project configuration
│   └── docker/           # Docker configurations
├── wp-content/           # WordPress content directory
│   ├── themes/          # Custom themes
│   ├── plugins/         # Custom plugins
│   └── uploads/         # Media uploads
└── wp-config.php        # WordPress configuration
```

## Configuration

wpld uses YAML configuration files stored in the `.wpld/` directory of each project. The main configuration includes:

- Database settings
- WordPress version
- PHP version
- Server configuration
- Network settings
- Container options

Example configuration:
```yaml
name: my-project
wordpress_version: latest
php_version: 8.1
database:
  name: wordpress
  user: wordpress
  password: wordpress
ports:
  web: 8080
  phpmyadmin: 8081
```

## Services

Each wpld project includes the following services:

- **WordPress**: The main WordPress application
- **MySQL**: Database server
- **PHPMyAdmin**: Web-based database management
- **Nginx**: Web server (if configured)

## Development Workflow

### Typical Development Session

```bash
# Create or start your project
wpld start

# Install dependencies
wpld wp plugin install --dev-dependencies

# Work on your code...

# View logs if needed
wpld logs wordpress

# Execute WP-CLI commands as needed
wpld wp db export backup-$(date +%Y%m%d).sql

# Stop when finished
wpld stop
```

### Database Management

Access PHPMyAdmin at `http://localhost:8081` (or your configured port) to manage your database graphically, or use WP-CLI:

```bash
# Export database
wpld wp db export backup.sql

# Import database
wpld wp db import backup.sql

# Reset database
wpld wp db reset --yes
```

## Troubleshooting

### Common Issues

1. **Port Conflicts**: If ports are already in use, wpld will prompt for alternative ports during setup.

2. **Docker Permission Issues**: Ensure Docker is running and your user has proper permissions.

3. **Container Startup Issues**: Check logs with `wpld logs [service]` for detailed error information.

### Debug Mode

Enable verbose logging for troubleshooting:

```bash
wpld -vv start  # Double verbose for maximum detail
```

### Clean Reset

If you encounter persistent issues:

```bash
# Stop all containers
wpld stop

# Remove containers and start fresh
wpld start --pull
```

## Contributing

We welcome contributions! Please feel free to submit issues, feature requests, or pull requests.

### Development Setup

```bash
# Clone and setup
git clone https://github.com/your-username/wpld.git
cd wpld

# Install dependencies
go mod download

# Build and test
make build
make test
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

### v1.0.0
- Initial release
- Core WordPress project management
- Docker container orchestration
- WP-CLI integration
- PHPMyAdmin support
- Network management

## Support

For questions, issues, or contributions:

- 📫 Create an issue on GitHub
- 💬 Join our community discussions
- 📖 Check the documentation

---

**wpld** - Making WordPress development with Docker simple and efficient! 🚀