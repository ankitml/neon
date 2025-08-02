Project Brief: Themed Quote Collections
This document outlines the project goals and implementation plan for creating "Themed Quote Collections," a web application designed to allow users to discover quotes based on various themes. The primary technical objective is to build and learn a powerful full-text search system using pg_search on a serverless PostgreSQL database.

1. Project Goals & Vision
Primary Goal: To build a functional full-stack application that implements efficient, real-world full-text search.

User Experience: Create a clean, fast, and intuitive interface where users can instantly find quotes related to themes like "courage," "innovation," or "perseverance."

Learning Objectives:

Master the fundamentals of PostgreSQL's full-text search capabilities (tsvector, to_tsquery, GIN indexes).

Gain practical experience building a robust backend API in Go.

Develop a modern, reactive frontend using SvelteKit and the Skeleton UI component library.

Understand the workflow of connecting a frontend application to a backend API, including handling CORS.

2. Core Features
A theme-based search bar for user input.

A dynamic results view that displays quotes in a clean, card-based layout.

Visual feedback for loading states and "no results found" scenarios.

A "Copy to Clipboard" button for each quote to allow for easy sharing.

A fully responsive design that works on both desktop and mobile devices.

3. Tech Stack
Backend: Go

Frontend: SvelteKit

UI Library: Skeleton UI (with Tailwind CSS)

Database: Neon (Serverless PostgreSQL with pg_search)

4. Implementation Plan
Phase 1: Backend Foundation & Database Setup (2-3 hours)
Goal: Create a Go server that can connect to the database and serve a basic health check endpoint.

Steps:

Setup Neon Database: Create a project, get the connection string, and create the quotes table. Import a dataset.

Initialize Go Project: Set up a backend folder and initialize a Go module.

Create Basic HTTP Server: Use net/http and the chi router to create a /health endpoint.

Establish Database Connection: Use the pgx driver to connect to Neon, pulling credentials from environment variables.

Phase 2: Implementing Search Logic (Go) (3-4 hours)
Goal: Create a search API endpoint that uses pg_search to find quotes.

Steps:

Prepare Database for Full-Text Search: Add a tsvector column (tsv), create an update trigger, and build a GIN index on the tsv column for performance.

Build the Search Endpoint: Create a GET /api/search route that accepts a q parameter and executes a to_tsquery SQL command.

Enable CORS: Add a CORS middleware to the chi router to allow requests from the SvelteKit frontend.

Phase 3: Frontend Scaffolding (SvelteKit & Skeleton UI) (2-3 hours)
Goal: Set up a new SvelteKit project with Skeleton UI and create the main page layout.

Steps:

Initialize SvelteKit Project: Run npm create svelte@latest frontend and choose the "Skeleton project" template.

Integrate Skeleton UI: Follow the official guide to install dependencies, set up Tailwind, and choose a theme.

Create Layout & Main Page: Build the main app shell in +layout.svelte and use Skeleton components in +page.svelte to create the search bar and results area.

Phase 4: Connecting Frontend to Backend (3-4 hours)
Goal: Make the SvelteKit app fetch and display search results from the Go API.

Steps:

Fetch Data: Write a function in +page.svelte that uses the fetch API to call the Go backend's /api/search endpoint.

Manage State: Use Svelte writable stores for searchResults, isLoading, and error.

Display Results: Use Svelte's {#each} and {#if} blocks to render the results using Skeleton's <Card> and <ProgressBar> components.

Phase 5: Polish & Refine (2+ hours)
Goal: Improve the user experience and add finishing touches.

Steps:

Refine UI: Add a header/footer, implement helpful user messages ("no results found"), and fine-tune typography.

Add "Copy to Clipboard" Feature: Add a copy button to each card and use a Skeleton <Toast> for user feedback.

Deployment Prep: Create a Dockerfile for the Go backend and ensure the SvelteKit frontend can be built for static deployment.
