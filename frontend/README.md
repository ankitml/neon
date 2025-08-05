# Quotes App

This is a modern web application for browsing and searching for quotes, built with Next.js and shadcn/ui.

## Architecture

The application is a Next.js-based frontend that communicates with a separate backend service to fetch quote data.

### Tech Stack

*   **Framework**: [Next.js](https://nextjs.org/) (using the App Router)
*   **Language**: [TypeScript](https://www.typescriptlang.org/)
*   **UI Components**: [shadcn/ui](https://ui.shadcn.com/) - A collection of beautifully designed, accessible, and customizable components built on top of Radix UI.
*   **Styling**: [Tailwind CSS](https://tailwindcss.com/) - A utility-first CSS framework for rapid UI development.
*   **Package Manager**: [pnpm](https://pnpm.io/)

### Project Structure

*   `app/`: Contains the application's pages and routes, following the Next.js App Router conventions.
*   `components/`: Contains the reusable UI components, with a dedicated `ui/` subdirectory for shadcn/ui components.
*   `lib/`: Contains utility functions and helper scripts.
*   `public/`: Stores static assets like images and fonts.

## Getting Started

Follow these instructions to get the development environment up and running.

### Prerequisites

*   Node.js (v18 or later)
*   [pnpm](https://pnpm.io/installation) package manager
*   A running instance of the backend service from the `backend/golang` directory, accessible at `http://localhost:8080`.

### Installation

1.  Navigate to the `quotes-app` directory:
    ```bash
    cd quotes-app
    ```

2.  Install the dependencies using pnpm:
    ```bash
    pnpm install
    ```

### Using the Makefile

This project includes a `Makefile` to simplify common development tasks. Here are the available commands:

*   `make install`: Install dependencies.
*   `make dev`: Start the development server in the background.
*   `make stop`: Stop the development server.
*   `make logs`: View the development server logs.
*   `make build`: Build the application for production.
*   `make start`: Start the production server.
*   `make lint`: Lint the codebase.
*   `make clean`: Remove `node_modules`, `.next`, and `dev-server.log`.

For example, to install the dependencies and start the development server, you can run:

```bash
make install
make dev
```

### Running the Development Server

1.  Start the Next.js development server in the background:
    ```bash
    make dev
    ```

2.  To view the server logs, run:
    ```bash
    make logs
    ```

3.  To stop the server, run:
    ```bash
    make stop
    ```

4.  Open your browser and navigate to [http://localhost:3000](http://localhost:3000) to see the application. The main quotes interface is available at `/quotes`.
