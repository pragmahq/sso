# Pragma SSO

A single sign-on solution tailored for Pragma, streamlining authentication for contributors across multiple pragma services.

## Development Setup

### Prerequisites

- [direnv](https://direnv.net) - For managing environment variables
- [Go](https://go.dev) (1.x or later) - Our primary programming language
- [Air](https://github.com/cosmtrek/air) - For live reloading during development
- PostgreSQL - Our database of choice

### Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/pragmagq/sso.git pragma-sso
   cd pragma-sso
   ```

2. Set up your environment:
   ```
   cp .envrc.example .envrc
   ```
   Edit `.envrc` with your local configuration.

3. Allow direnv to load the environment:
   ```
   direnv allow .
   ```

4. Start the development server:
   ```
   air -c .air.toml
   ```

This will launch the application in development mode with live reloading.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). See the [LICENSE](LICENSE) file for details.

## Contributing

We welcome contributions to Pragma SSO! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Please make sure to update tests as appropriate and adhere to our coding standards.

For major changes, please open an issue first to discuss what you would like to change. Ensure to update the README.md with details of changes to the interface, if applicable.

## Support

If you encounter any problems or have any questions, please open an issue in the GitHub repository.
