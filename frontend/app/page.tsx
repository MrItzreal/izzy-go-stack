import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import Link from "next/link"
import { UserButton, SignedIn, SignedOut } from "@clerk/nextjs"
import { Github } from "lucide-react"

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-4 md:p-24">
      <div className="z-10 w-full max-w-5xl items-center justify-between font-mono text-sm flex">
        <p className="fixed left-0 top-0 flex w-full justify-center border-b border-gray-300 bg-gradient-to-b from-zinc-200 pb-6 pt-8 backdrop-blur-2xl dark:border-neutral-800 dark:bg-zinc-800/30 dark:from-inherit lg:static lg:w-auto lg:rounded-xl lg:border lg:bg-gray-200 lg:p-4 lg:dark:bg-zinc-800/30">
          Full-Stack Boilerplate
        </p>
        <div className="fixed right-0 top-0 flex items-center justify-center gap-2 p-4">
          <SignedIn>
            <UserButton afterSignOutUrl="/" />
          </SignedIn>
          <SignedOut>
            <Link href="/sign-in">
              <Button>Sign In</Button>
            </Link>
          </SignedOut>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-12 w-full max-w-5xl">
        <Card>
          <CardHeader>
            <CardTitle>Frontend Stack</CardTitle>
            <CardDescription>Modern, type-safe frontend technologies</CardDescription>
          </CardHeader>
          <CardContent>
            <ul className="list-disc pl-5 space-y-2">
              <li>TypeScript for type safety</li>
              <li>React for component-based UI</li>
              <li>Next.js for server-side rendering</li>
              <li>Tailwind CSS for styling</li>
              <li>Clerk for authentication</li>
              <li>Zod for validation</li>
            </ul>
          </CardContent>
          <CardFooter>
            <Link href="https://github.com/your-username/your-repo" target="_blank">
              <Button variant="outline" className="gap-2">
                <Github size={16} />
                View on GitHub
              </Button>
            </Link>
          </CardFooter>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Backend Stack</CardTitle>
            <CardDescription>Performant and scalable backend technologies</CardDescription>
          </CardHeader>
          <CardContent>
            <ul className="list-disc pl-5 space-y-2">
              <li>Golang for API server</li>
              <li>Supabase (PostgreSQL) for database</li>
              <li>Drizzle ORM for database interactions</li>
              <li>Zod (Go port) for validation</li>
              <li>Stripe for payment processing</li>
            </ul>
          </CardContent>
          <CardFooter>
            <Link href="/api-docs">
              <Button variant="outline">View API Docs</Button>
            </Link>
          </CardFooter>
        </Card>
      </div>
    </main>
  )
}
