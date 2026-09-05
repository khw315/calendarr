import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cn } from "@/utils"

export interface ButtonProps
    extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    asChild?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
    ({ className, asChild = false, ...props }, ref) => {
        const Comp = asChild ? Slot : "button"
        return (
            <Comp
                className={cn(
                    "inline-flex items-center justify-center whitespace-nowrap rounded-base text-sm font-base ring-offset-white transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-black focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
                    "border-2 border-border bg-main text-main-foreground shadow-shadow hover:translate-x-[boxShadowX] hover:translate-y-[boxShadowY] hover:shadow-none active:translate-x-[boxShadowX] active:translate-y-[boxShadowY] active:shadow-none",
                    "h-10 px-4 py-2 uppercase font-heading",
                    className
                )}
                ref={ref}
                {...props}
            />
        )
    }
)
Button.displayName = "Button"

export { Button }
