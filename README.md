# GoLang Web Scraper

This is a GoLang application that allows you to scrape web content from one or more URLs. It retrieves the HTML content of the specified URLs and saves it to local files. Additionally, it can extract metadata from the HTML, such as the number of links and images present on the page.

## Prerequisites

Before running this application, make sure you have the following installed:

- GoLang (version 1.16 or higher)

## Usage

To use this application, follow the steps below:

1. Clone the repository or copy the code into a local file.

1. Open a terminal or command prompt and navigate to the directory containing the Go file.

1. Build the application by running the following command:

   ```bash
   go build
   ```

1. Run the application with the desired command-line arguments. The supported options are:

    - `--metadata`: Enables metadata mode, which extracts additional information from the HTML content.

    - `<URLs>`: Provide one or more URLs as command-line arguments, separated by spaces.

   Example usage:

   ```bash
   ./web-scraper --metadata https://www.example.com https://www.another-example.com
   ```

   Replace `web-scraper` with the name of the built executable file.

1. The application will retrieve the HTML content of the specified URLs, save it to local files, and display metadata (if enabled). The HTML content will be saved as `<hostname>.html`, and the associated resources (images, stylesheets, etc.) will be saved in a folder named `<hostname>_content`.

## Features

- Fetches HTML content from one or more URLs.

- Saves HTML content to local files.

- Extracts metadata from HTML, including the number of links and images.

- Automatically downloads associated resources (images, stylesheets, etc.) and updates the HTML with the correct local URLs.

## Customization

You can customize the behavior of the application by modifying the Go code. For example, you can add additional metadata extraction logic or enhance the resource downloading process.

## License

This application is open-source and distributed under the [MIT License](https://opensource.org/licenses/MIT).

## Disclaimer

This application is provided as-is without any warranty. Use it at your own risk.