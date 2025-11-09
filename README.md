# WorkZen HRMS - FrontendThis is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

Modern React-based frontend for WorkZen HRMS built with Next.js 16, TypeScript, and Shadcn UI.## Getting Started

## 🎨 Tech StackFirst, run the development server:

- **Next.js 16** - React framework with Turbopack```bash

- **TypeScript** - Type-safe JavaScriptnpm run dev

- **Shadcn UI** - Radix UI component library# or

- **Recharts 2.15.4** - Data visualizationyarn dev

- **Tailwind CSS** - Utility-first CSS# or

- **Tabler Icons** - Icon librarypnpm dev

- **Sonner** - Toast notifications# or

bun dev

## 🚀 Quick Start```

### PrerequisitesOpen [http://localhost:3000](http://localhost:3000) with your browser to see the result.

- Node.js 18.x or higher

- pnpm (recommended) or npmYou can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

### InstallationThis project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

```bash## Learn More

# Install dependencies

pnpm installTo learn more about Next.js, take a look at the following resources:



# Run development server- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.

pnpm dev- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.



# Build for productionYou can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

pnpm build

## Deploy on Vercel

# Start production server

pnpm startThe easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

```

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.

The app will run on [http://localhost:3000](http://localhost:3000)

## 📁 Project Structure

```
frontend/
├── app/                    # Next.js App Router
│   ├── dashboard/         # Protected dashboard routes
│   │   ├── users/         # User management
│   │   ├── attendance/    # Attendance tracking
│   │   ├── leaves/        # Leave management
│   │   ├── payroll/       # Payroll processing
│   │   ├── departments/   # Department management
│   │   ├── documents/     # Document management
│   │   ├── profile/       # User profile
│   │   └── settings/      # Settings
│   ├── login/             # Login page
│   ├── signup/            # Company registration
│   ├── verify-email/      # Email verification
│   └── layout.tsx         # Root layout
├── components/
│   ├── ui/                # Shadcn UI components
│   ├── dashboard/         # Dashboard components
│   ├── app-sidebar.tsx    # Main sidebar
│   ├── login-form.tsx     # Login form
│   └── signup-form.tsx    # Signup form
├── lib/
│   ├── api-service.ts     # API client with auth
│   ├── auth-context.tsx   # Auth state management
│   ├── config.ts          # API endpoints
│   ├── types.ts           # TypeScript types
│   └── utils.ts           # Utility functions
└── hooks/                  # Custom React hooks
```

## 🔑 Key Features

### Authentication

- JWT-based authentication
- Email verification
- Protected routes with middleware
- Auto-redirect on unauthorized access

### Dashboard

- Interactive charts with Recharts
- Real-time statistics
- Department-wise analytics
- Responsive 3-column layout
- Role-based data filtering

### UI Components

- Consistent table styling across all pages
- Standardized form layouts
- Loading states and skeletons
- Toast notifications for feedback
- Error boundaries for error handling
- Dark mode support

### Data Management

- CRUD operations for all entities
- File upload with preview
- Advanced filtering and search
- Pagination support
- Real-time updates

## 🔧 Configuration

### Environment Variables

Create `.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:5000
```

### API Integration

All API calls go through `lib/api-service.ts`:

```typescript
import { apiService } from "@/lib/api-service";

// Get data
const users = await apiService.get("/api/v1/users");

// Post data
const result = await apiService.post("/api/v1/users", userData);
```

## 📦 Available Scripts

```bash
# Development
pnpm dev          # Start dev server with Turbopack
pnpm build        # Build for production
pnpm start        # Start production server

# Linting & Formatting
pnpm lint         # Run ESLint
pnpm type-check   # Run TypeScript compiler

# Testing
pnpm test         # Run tests (if configured)
```

## 🎨 Styling Guidelines

### Using Tailwind CSS

```tsx
<div className="flex items-center gap-4 p-6 rounded-lg border bg-card">
  {/* Content */}
</div>
```

### Using Shadcn Components

```tsx
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";

<Card>
  <CardHeader>Title</CardHeader>
  <CardContent>
    <Button variant="default" size="sm">
      Click Me
    </Button>
  </CardContent>
</Card>;
```

### Responsive Design

```tsx
// Mobile-first approach
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
  {/* Responsive grid */}
</div>
```

## 🔐 Authentication Flow

1. User logs in → JWT token stored in localStorage
2. Token included in all API requests via Authorization header
3. Middleware checks auth on protected routes
4. Auto-redirect to login on 401 Unauthorized
5. Token refresh on expiry (if implemented)

## 📊 Chart Components

Using Recharts for data visualization:

```tsx
import { BarChart, Bar, XAxis, YAxis, Tooltip, Legend } from "recharts";

<ResponsiveContainer width="100%" height={350}>
  <BarChart data={data}>
    <XAxis dataKey="name" />
    <YAxis />
    <Tooltip />
    <Legend />
    <Bar dataKey="value" fill="#3b82f6" />
  </BarChart>
</ResponsiveContainer>;
```

## 🚀 Deployment

### Vercel (Recommended)

1. Push code to GitHub
2. Import repository in Vercel
3. Set environment variables
4. Deploy

### Self-Hosted

```bash
# Build the app
pnpm build

# Start production server
pnpm start

# Or use PM2
pm2 start npm --name "workzen-frontend" -- start
```

## 📝 Adding New Pages

1. Create page in `app/dashboard/[page-name]/page.tsx`
2. Add route to sidebar in `components/app-sidebar.tsx`
3. Create necessary components in `components/`
4. Add types in `lib/types.ts`
5. Update API service if needed

## 🐛 Debugging

### Common Issues

**Build Errors:**

```bash
# Clear Next.js cache
rm -rf .next
pnpm build
```

**Type Errors:**

```bash
# Check TypeScript
pnpm type-check
```

**Module Not Found:**

```bash
# Reinstall dependencies
rm -rf node_modules pnpm-lock.yaml
pnpm install
```

## 📚 Learn More

- [Next.js Documentation](https://nextjs.org/docs)
- [Shadcn UI](https://ui.shadcn.com/)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [Recharts](https://recharts.org/)
- [TypeScript](https://www.typescriptlang.org/docs)

## 🤝 Contributing

See the main project [README](../README.md) for contribution guidelines.

---

**Part of WorkZen HRMS** | [Backend](../backend) | [Main README](../README.md)
