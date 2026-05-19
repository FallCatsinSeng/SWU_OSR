import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

const badgeVariants = cva(
  "inline-flex items-center rounded-geist-full px-2 py-0.5 text-xs font-normal transition-colors",
  {
    variants: {
      variant: {
        // Default — soft canvas background (badge-secondary from Geist)
        default: "bg-geist-canvas-soft text-geist-body",
        // Secondary — slightly deeper
        secondary: "bg-geist-canvas-soft-2 text-geist-body",
        // Destructive
        destructive: "bg-geist-error-soft text-geist-error-deep",
        // Outline — hairline border
        outline: "border border-geist-hairline text-geist-body",
        // Link / info
        info: "bg-geist-link-bg-soft text-geist-link-deep",
        // Warning
        warning: "bg-geist-warning-soft text-geist-warning-deep",
        // Success
        success: "bg-emerald-50 text-emerald-700",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
);

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <div className={cn(badgeVariants({ variant }), className)} {...props} />
  );
}

export { Badge, badgeVariants };
