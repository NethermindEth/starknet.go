import {
  RiGithubFill,
  RiSlackFill,
  RiTelegram2Fill,
  RiTwitterXFill,
  RiYoutubeFill,
} from "@remixicon/react"
import { StarknetGoLogo } from "../../public/StarknetGoLogo"
const CURRENT_YEAR = new Date().getFullYear()

const Footer = () => {
  const sections = {
    solutions: {
      title: "Solutions",
      items: [
        { label: "Starknet.go", href: "#" },
        { label: "Voyager", href: "https://voyager.online/" },
        { label: "Juno", href: "https://www.nethermind.io/juno" },
        { label: "CairoVM - Go", href: "https://github.com/NethermindEth/cairo-vm-go" },
        { label: "Starknet RPC", href: "https://data.voyager.online/"},
        { label: "Starkweb", href: "https://www.starkweb.xyz/" },
      ],
    },
    company: {
      title: "Company",
      items: [
        { label: "About Nethermind", href: "https://www.nethermind.io/" },
        { label: "Blog", href: "https://www.nethermind.io/blog" },
        { label: "Careers", href: "https://www.nethermind.io/open-roles" },
        { label: "Events", href: "https://www.nethermind.io/events" },
      ],
    },
    resources: {
      title: "Resources",
      items: [
        {
          label: "Community",
          href: "#",
          external: true,
        },
        { label: "Contact Us", href: "https://www.nethermind.io/contact-us#" },
        { label: "Community", href: "https://discord.com/invite/PaCMRFdvWT" },
        { label: "Legal", href: "https://www.nethermind.io/legal" },
        { label: "Media Kit", href: "https://drive.google.com/drive/folders/1pGJw5TAjo8M1RdGVbqPrEjmvkpwalIfL" },
      ],
    },
    ecosystem: {
      title: "Ecosystem",
      items: [
        { label: "Starknet", href: "https://www.starknet.io/", external: true },
        { label: "Github", href: "https://github.com/NethermindEth/starknet.go", external: true },
        { label: "NethermindStark", href: "https://x.com/NethermindStark?ref_src=twsrc%5Egoogle%7Ctwcamp%5Eserp%7Ctwgr%5Eauthor", external: true },
        { label: "Telegram", href: "https://t.me/StarknetGo", external: true },
      ],
    },
  }

  return (
    <div className="px-4 xl:px-0">
      <footer
        id="footer"
        className="relative mx-auto flex max-w-6xl flex-wrap pt-4"
      >
        {/* Vertical Lines */}
        <div className="pointer-events-none inset-0">
          {/* Left */}
          <div
            className="absolute inset-y-0 my-[-5rem] w-px"
            style={{
              maskImage: "linear-gradient(transparent, white 5rem)",
            }}
          >
            <svg className="h-full w-full" preserveAspectRatio="none">
              <line
                x1="0"
                y1="0"
                x2="0"
                y2="100%"
                className="stroke-gray-300"
                strokeWidth="2"
                strokeDasharray="3 3"
              />
            </svg>
          </div>

          {/* Right */}
          <div
            className="absolute inset-y-0 right-0 my-[-5rem] w-px"
            style={{
              maskImage: "linear-gradient(transparent, white 5rem)",
            }}
          >
            <svg className="h-full w-full" preserveAspectRatio="none">
              <line
                x1="0"
                y1="0"
                x2="0"
                y2="100%"
                className="stroke-gray-300"
                strokeWidth="2"
                strokeDasharray="3 3"
              />
            </svg>
          </div>
        </div>
        <svg
          className="mb-10 h-20 w-full border-y border-dashed border-gray-300 stroke-gray-300"
          // style={{
          //   maskImage:
          //     "linear-gradient(transparent, white 10rem, white calc(100% - 10rem), transparent)",
          // }}
        >
          <defs>
            <pattern
              id="diagonal-footer-pattern"
              patternUnits="userSpaceOnUse"
              width="64"
              height="64"
            >
              {Array.from({ length: 17 }, (_, i) => {
                const offset = i * 8
                return (
                  <path
                    key={i}
                    d={`M${-106 + offset} 110L${22 + offset} -18`}
                    stroke=""
                    strokeWidth="1"
                  />
                )
              })}
            </pattern>
          </defs>
          <rect
            stroke="none"
            width="100%"
            height="100%"
            fill="url(#diagonal-footer-pattern)"
          />
        </svg>
        <div className="mr-auto flex w-full justify-between lg:w-fit lg:flex-col">
          <a
            href="/"
            className="flex items-center font-medium text-gray-700 select-none sm:text-sm"
          >
            <StarknetGoLogo className="ml-2 w-48" />

            <span className="sr-only">Starknet.go Logo (go home)</span>
          </a>

          <div>
            <div className="mt-4 flex items-center">
              {/* Social Icons */}
              <a
                href="https://x.com/NethermindStark?ref_src=twsrc%5Egoogle%7Ctwcamp%5Eserp%7Ctwgr%5Eauthor"
                target="_blank"
                rel="noopener noreferrer"
                className="rounded-sm p-2 text-gray-700 transition-colors duration-200 hover:bg-gray-200 hover:text-gray-900"
              >
                <RiTwitterXFill className="size-5" />
                  </a>
              <a
                href="https://www.youtube.com/channel/UCgIPcx1C29j8IUx7BtF77YQ"
                target="_blank"
                rel="noopener noreferrer"
                className="rounded-sm p-2 text-gray-700 transition-colors duration-200 hover:bg-gray-200 hover:text-gray-900"
              >
                <RiYoutubeFill className="size-5" />
              </a>
              <a
                href="https://github.com/NethermindEth/starknet.go#"
                target="_blank"
                rel="noopener noreferrer"
                className="rounded-sm p-2 text-gray-700 transition-colors duration-200 hover:bg-gray-200 hover:text-gray-900"
              >
                <RiGithubFill className="size-5" />
              </a>
              <a
                href="https://t.me/StarknetGo"
                target="_blank"
                rel="noopener noreferrer"
                className="rounded-sm p-2 text-gray-700 transition-colors duration-200 hover:bg-gray-200 hover:text-gray-900"
              >
              <RiTelegram2Fill className="size-5" />
              </a>
            </div>
            <div className="ml-2 hidden text-sm text-gray-700 lg:inline">
              &copy; {CURRENT_YEAR} Nethermind
            </div>
          </div>
        </div>

        {/* Footer Sections */}
        {Object.entries(sections).map(([key, section]) => (
          <div key={key} className="mt-10 min-w-44 pl-2 lg:mt-0 lg:pl-0">
            <h3 className="mb-4 font-medium text-gray-900 sm:text-sm">
              {section.title}
            </h3>
            <ul className="space-y-4">
              {section.items.map((item) => (
                <li key={item.label} className="text-sm">
                  <a
                    href={item.href}
                    className="text-gray-600 transition-colors duration-200 hover:text-gray-900"
                  >
                    {item.label}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </footer>
    </div>
  )
}

export default Footer
