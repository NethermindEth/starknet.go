// Tremor Divider [v0.0.2]

import React from "react"

import { cx } from "../lib/utils"

type DividerProps = React.ComponentPropsWithoutRef<"div">

const Divider = React.forwardRef<HTMLDivElement, DividerProps>(
  ({ className, children, ...props }, forwardedRef) => (
    <div
      ref={forwardedRef}
      className={cx(
        // base
        "mx-auto my-6 flex w-full items-center justify-between gap-3 text-sm",
        // text color
        "text-gray-500",
        className,
      )}
      tremor-id="tremor-raw"
      {...props}
    >
      {children ? (
        <>
          <div
            className={cx(
              // base
              "h-[1px] w-full",
              // background color
              "bg-linear-to-r from-transparent to-gray-200",
            )}
          />
          <div className="whitespace-nowrap text-inherit">{children}</div>
          <div
            className={cx(
              // base
              "h-[1px] w-full",
              // background color
              "bg-linear-to-l from-transparent to-gray-200",
            )}
          />
        </>
      ) : (
        <div
          className={cx(
            // base
            "h-[1px] w-full",
            // background color
            "bg-linear-to-l from-transparent via-gray-200 to-transparent",
          )}
        />
      )}
    </div>
  ),
)

Divider.displayName = "Divider"

export { Divider }
