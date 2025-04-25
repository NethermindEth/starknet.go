import {
  RiCodeSSlashFill,
  RiLockFill,
  RiSpeedFill,
  RiStackFill,
} from "@remixicon/react"
import { Divider } from "../Divider"
import AnalyticsIllustration from "./AnalyticsIllustration"
import { StickerCard } from "./StickerCard"

export function SolarAnalytics() {
  return (
    <section
      aria-labelledby="starknet-analytics"
      className="relative mx-auto w-full max-w-6xl overflow-hidden"
    >
      <div>
        <h2
          id="starknet-analytics"
          className="relative scroll-my-24 text-lg font-semibold tracking-tight text-orange-500"
        >
          Starknet Analytics
          <div className="absolute top-1 -left-[8px] h-5 w-[3px] rounded-r-sm bg-orange-500" />
        </h2>
        <p className="mt-2 max-w-lg text-3xl font-semibold tracking-tighter text-balance text-gray-900 md:text-4xl">
          Transform blockchain data into actionable insights with Go-powered tools
        </p>
      </div>
      <div className="*:pointer-events-none">
        <AnalyticsIllustration />
      </div>
      <Divider className="mt-0"></Divider>
      <div className="grid grid-cols-1 grid-rows-2 gap-6 md:grid-cols-4 md:grid-rows-1">
        <StickerCard
          Icon={RiSpeedFill}
          title="High Performance"
          description="Native Go implementation delivering exceptional speed and efficiency."
        />
        <StickerCard
          Icon={RiStackFill}
          title="Scalable Architecture"
          description="Built to handle high transaction volumes with minimal resource usage."
        />
        <StickerCard
          Icon={RiLockFill}
          title="Enhanced Security"
          description="Type-safe implementation with comprehensive security features."
        />
        <StickerCard
          Icon={RiCodeSSlashFill}
          title="Developer Friendly"
          description="Intuitive APIs and extensive documentation for rapid development."
        />
      </div>
    </section>
  )
}
