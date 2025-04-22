# Full-Stack Boilerplate

This is a boilerplate project that provides a solid foundation for building modern web applications with a robust tech stack.

## Tech Stack

### Frontend

- **TypeScript**: For type safety and improved code maintainability.
- **React**: For building a component-based, interactive user interface.
- **Next.js**: For server-side rendering, routing, and API routes.
- **Tailwind CSS**: For rapid UI development and responsive styling.
- **Clerk**: For user authentication and management.
- **Zod**: For runtime data validation and schema definition.

### Backend

- **Golang**: For building a performant and scalable API server.
- **Supabase (PostgreSQL)**: As the database, providing a managed PostgreSQL instance with features like authentication and storage.
- **Drizzle ORM**: For type-safe database interactions.
- **Zod (Go port or similar)**: For validating incoming data from the frontend.
- **Stripe**: For handling payment processing.

## Project Structure

\`\`\`
├── frontend/               # Next.js frontend application
│   ├── app/                # App Router pages and layouts
│   ├── components/         # React components
│   ├── lib/                # Utility functions and shared code
│   ├── public/             # Static assets
│   └── ...
├── backend/                # Golang API server
│   ├── cmd/                # Application entry points
│   ├── internal/           # Internal packages
│   ├── pkg/                # Reusable packages
│   └── ...
└── README.md               # Project documentation
\`\`\`

## Getting Started

### Prerequisites

- Node.js 18+ and npm/yarn/pnpm
- Go 1.20+
- Supabase account
- Clerk account
- Stripe account

### Setup Instructions

1. Clone this repository
2. Set up environment variables (see `.env.example` files in both frontend and backend directories)
3. Install frontend dependencies:
   \`\`\`bash
   cd frontend
   npm install
   \`\`\`
4. Install backend dependencies:
   \`\`\`bash
   cd backend
   go mod tidy
   \`\`\`
5. Start the development servers:
   - Frontend: `npm run dev` in the frontend directory
   - Backend: `go run cmd/api/main.go` in the backend directory

## Features

- User authentication with Clerk
- Database interactions with Drizzle ORM
- Type-safe API requests with Zod validation
- Payment processing with Stripe
- Responsive UI with Tailwind CSS

## License

MIT
