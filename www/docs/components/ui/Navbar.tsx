"use client"

import { siteConfig } from "@/app/siteConfig"
import useScroll from "../lib/useScroll"
import { cx } from "../lib/utils"
import { RiCloseFill, RiMenuFill } from "@remixicon/react"
import Link from "next/link"
import React from "react"
import { SolarLogo } from "../../../public/SolarLogo"
import { Button } from "../Button"

export function NavBar() {
  const [open, setOpen] = React.useState(false)
  const scrolled = useScroll(15)

  return (
    <header
      className={cx(
        "fixed inset-x-4 top-4 z-50 mx-auto flex max-w-6xl justify-center rounded-lg border border-transparent px-3 py-3 transition duration-300",
        scrolled || open
          ? "border-gray-200/50 bg-white/80 shadow-2xl shadow-black/5 backdrop-blur-sm"
          : "bg-white/0",
      )}
    >
      <div className="w-full md:my-auto">
        <div className="relative flex items-center justify-between">
          <Link href={siteConfig.baseLinks.home} aria-label="Home">
            <span className="sr-only">Solar Tech Logo</span>
            <SolarLogo className="w-22" />
          </Link>
          <nav className="hidden sm:block md:absolute md:top-1/2 md:left-1/2 md:-translate-x-1/2 md:-translate-y-1/2 md:transform">
            <div className="flex items-center gap-10 font-medium">
              <Link className="px-2 py-1 text-gray-900" href="#solutions">
                Solutions
              </Link>
              <Link className="px-2 py-1 text-gray-900" href="#farm-management">
                Farm Management
              </Link>
              <Link className="px-2 py-1 text-gray-900" href="#solar-analytics">
                Analytics
              </Link>
            </div>
          </nav>
          <Button
            variant="secondary"
            className="hidden h-10 font-semibold sm:block"
          >
            Get a quote
          </Button>
          <Button
            onClick={() => setOpen(!open)}
            variant="secondary"
            className="p-1.5 sm:hidden"
            aria-label={open ? "CloseNavigation Menu" : "Open Navigation Menu"}
          >
            {!open ? (
              <RiMenuFill
                className="size-6 shrink-0 text-gray-900"
                aria-hidden
              />
            ) : (
              <RiCloseFill
                className="size-6 shrink-0 text-gray-900"
                aria-hidden
              />
            )}
          </Button>
        </div>
        <nav
          className={cx(
            "mt-6 flex flex-col gap-6 text-lg ease-in-out will-change-transform sm:hidden",
            open ? "" : "hidden",
          )}
        >
          <ul className="space-y-4 font-medium">
            <li onClick={() => setOpen(false)}>
              <Link href="#solutions">Solutions</Link>
            </li>
            <li onClick={() => setOpen(false)}>
              <Link href="#farm-management">Farm Management</Link>
            </li>
            <li onClick={() => setOpen(false)}>
              <Link href="#solar-analytics">Analytics</Link>
            </li>
          </ul>
          <Button variant="secondary" className="text-lg">
            Get a quote
          </Button>
        </nav>
      </div>
    </header>
  )
}
