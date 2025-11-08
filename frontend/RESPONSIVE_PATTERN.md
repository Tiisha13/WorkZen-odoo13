# Responsive Dialog/Drawer Pattern

## Overview

The user management page now implements a responsive UI pattern that provides optimal user experience across different device sizes.

## Pattern Components

### 1. **useMediaQuery Hook** (`lib/hooks.ts`)

A custom React hook that detects the current viewport size:

```typescript
const isMobile = useMediaQuery("(max-width: 768px)");
```

- Returns `true` when viewport width is ≤ 768px (mobile/tablet)
- Returns `false` for desktop screens
- Automatically updates when window is resized
- SSR-safe with proper initialization

### 2. **Conditional Rendering**

Based on the `isMobile` value, we render different components:

#### Mobile (Drawer)

- Slides up from bottom of screen
- Better for touch interactions
- Easier thumb access on mobile devices
- Natural swipe-down to dismiss

#### Desktop (Dialog)

- Centered modal overlay
- Better for mouse/keyboard interaction
- More screen real estate for form fields
- Traditional web application UX

### 3. **Alert Dialog for Confirmations**

For destructive actions (like delete), we always use AlertDialog:

- Clear warning message
- Shows user details being deleted
- Cancel/Confirm buttons
- Prevents accidental deletions
- Same behavior on all devices

## Implementation in User Management

### State Management

```typescript
const [isDialogOpen, setIsDialogOpen] = useState(false);
const [deleteUser, setDeleteUser] = useState<User | null>(null);
const isMobile = useMediaQuery("(max-width: 768px)");
```

### Form Component

The form fields are extracted into a reusable component:

```typescript
const UserFormFields = () => (
  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
    {/* All form fields */}
  </div>
);
```

This component is used in both Drawer and Dialog, ensuring consistency.

### Delete Flow

1. User clicks delete button → `setDeleteUser(user)`
2. AlertDialog opens showing user details
3. User confirms → `handleDelete(user.id)` called
4. User cancels → `setDeleteUser(null)` closes dialog

## Benefits

✅ **Better Mobile UX**: Drawer is easier to use on touch devices
✅ **Better Desktop UX**: Dialog provides traditional web experience
✅ **Code Reusability**: Single form component used in both contexts
✅ **Safety**: AlertDialog prevents accidental deletions
✅ **Consistency**: Same functionality across all screen sizes
✅ **Accessibility**: Both Dialog and Drawer are accessible components from Shadcn UI

## Applying to Other Pages

To apply this pattern to other CRUD pages (Departments, Documents, etc.):

1. **Add imports**:

```typescript
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useMediaQuery } from "@/lib/hooks";
```

2. **Add state**:

```typescript
const isMobile = useMediaQuery("(max-width: 768px)");
const [deleteItem, setDeleteItem] = useState<YourType | null>(null);
```

3. **Extract form fields** to a component

4. **Replace Dialog** with conditional rendering (Drawer for mobile, Dialog for desktop)

5. **Update delete handler** to use `setDeleteItem(item)` instead of window.confirm

6. **Add AlertDialog** for delete confirmations

## Example Usage

```typescript
// Open form for editing
setEditingUser(user);
setFormData({ ...user });
setIsDialogOpen(true);

// Trigger delete confirmation
setDeleteUser(user);

// Handle delete (called from AlertDialog)
const handleDelete = async (id: string) => {
  await apiService.delete(`${API_ENDPOINTS.USERS}/${id}`);
  setDeleteUser(null);
  fetchUsers();
};
```

## Testing Checklist

- [ ] Desktop: Dialog opens centered with proper styling
- [ ] Desktop: Form submission works
- [ ] Desktop: Cancel closes dialog
- [ ] Mobile: Drawer slides up from bottom
- [ ] Mobile: Form submission works
- [ ] Mobile: Swipe down or cancel closes drawer
- [ ] Delete: AlertDialog shows correct user details
- [ ] Delete: Cancel closes without deleting
- [ ] Delete: Confirm successfully deletes user
- [ ] Resize: Switches between Drawer/Dialog at 768px breakpoint

## Browser Support

- Modern browsers with matchMedia support
- Chrome, Firefox, Safari, Edge (recent versions)
- Mobile browsers (iOS Safari, Chrome Android)
