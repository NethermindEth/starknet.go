import { RiArrowRightUpLine } from "@remixicon/react"
import { FadeContainer, FadeDiv, FadeSpan } from "../Fade"
import GameOfLife from "./HeroBackground"

export function Hero() {
  return (
    <>
      <style>
        {`
          .scrollbar-hide {
            -ms-overflow-style: none;
            scrollbar-width: none;
          }
          .scrollbar-hide::-webkit-scrollbar {
            display: none;
          }
        `}
      </style>
      <section aria-label="hero" className="relative h-screen flex items-center justify-center -mt-46 bg-gradient-to-b from-gray-50 to-white">
      {/* Background Animation - positioned absolutely behind everything */}
      <div className="absolute inset-0 w-full h-full overflow-hidden">
        <GameOfLife />
      </div>
      
      {/* Main Content - positioned above the background */}
      <FadeContainer className="relative z-10 flex flex-col items-center justify-center">
        <FadeDiv className="mx-auto">
          <a
            aria-label="View latest update the changelog page"
            href="https://github.com/NethermindEth/starknet.go/releases/tag/v0.10.0"
            target="_blank"
            rel="noopener noreferrer"
            className="mx-auto w-full"
          >
            <div className="inline-flex max-w-full items-center gap-3 rounded-full bg-white/5 px-2.5 py-0.5 pr-3 pl-0.5 font-medium text-gray-900 ring-1 shadow-lg shadow-orange-400/20 ring-black/10 filter backdrop-blur-[1px] transition-colors hover:bg-orange-500/[2.5%] focus:outline-hidden sm:text-sm">
              <span className="shrink-0 truncate rounded-full border bg-gray-50 px-2.5 py-1 text-sm text-gray-600 sm:text-xs">
                Latest
              </span>
              <span className="flex items-center gap-1 truncate">
                <span className="w-full truncate">
                  Starknet.go v0.10.0 Released
                </span>

                <RiArrowRightUpLine className="size-4 shrink-0 text-gray-700" />
              </span>
            </div>
          </a>
        </FadeDiv>
        <h1 className="mt-8 text-center text-5xl font-semibold tracking-tighter text-gray-900 sm:text-8xl sm:leading-[5.5rem]">
          <FadeSpan>Starknet</FadeSpan> <FadeSpan>for</FadeSpan>
          <br />
          <FadeSpan>Go</FadeSpan> <FadeSpan>Developers</FadeSpan>
        </h1>
        <p className="mt-5 pt-4 max-w-xl text-center text-base text-balance text-gray-700 sm:mt-8 sm:text-xl">
          <FadeSpan>Building the future of Starknet with a powerful</FadeSpan>{" "}
          <FadeSpan>Go implementation for scalable and efficient</FadeSpan>{" "}
          <FadeSpan>blockchain development.</FadeSpan>
        </p>
        
        {/* Code snippet section */}
        <FadeDiv className="mt-8">
          <div className="relative group">
            <div className="absolute inset-0 bg-gradient-to-r from-[var(--stkgo-yellow)] to-[var(--stkgo-orange)] rounded-lg blur opacity-20 group-hover:opacity-30 transition duration-300"></div>
            <div className="relative bg-gray-900 rounded-lg px-4 py-3 border border-gray-200/10">
              <div className="flex items-center justify-between">
                <code className="text-green-400 font-mono text-sm sm:text-base">
                  go get github.com/NethermindEth/starknet.go
                </code>
                <button 
                  onClick={() => navigator.clipboard.writeText('go get github.com/NethermindEth/starknet.go')}
                  className="ml-3 text-gray-400 hover:text-white transition-colors duration-200 p-1 rounded"
                  title="Copy to clipboard"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </FadeDiv>
        
        <FadeDiv>
          <a
            className="mt-6 inline-flex cursor-pointer flex-row items-center justify-center gap-1 rounded-md border-b-[1.5px] border-[var(--stkgo-orange)] bg-linear-to-b from-[var(--stkgo-yellow)] to-[var(--stkgo-orange)] px-5 py-3 leading-4 font-medium tracking-wide whitespace-nowrap text-white shadow-[0_0_0_2px_rgba(0,0,0,0.04),0_0_14px_0_rgba(255,255,255,0.19)] transition-all duration-200 ease-in-out hover:shadow-orange-300"
            href="./docs/introduction/getting-started"
          >
            Get Started
          </a>
        </FadeDiv>
      </FadeContainer>
    </section>
    </>
  )
}