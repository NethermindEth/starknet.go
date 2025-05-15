import { Button } from "../Button"

export function CallToAction() {
  return (
    <section aria-labelledby="cta-title" className="mx-auto max-w-6xl">
      <div className="grid items-center gap-8 sm:grid-cols-6">
        <div className="sm:col-span-2">
          <h2
            id="cta-title"
            className="scroll-my-60 text-3xl font-semibold tracking-tighter text-balance text-gray-900 md:text-4xl"
          >
            Ready to build on Starknet?
          </h2>
          <p className="mt-3 mb-8 text-lg text-gray-600">
            Start developing with Starknet.go today or connect with our engineers
            to discuss your blockchain implementation needs.
          </p>
          <div className="flex flex-wrap gap-4 pt-10">
            <Button asChild className="text-md">
              <a href="/docs/introduction/getting-started">Get started</a>
            </Button>
            {/* <Button asChild className="text-md" variant="secondary">
              <a href="/docs/introduction/installation">View documentation</a>
            </Button> */}
          </div>
        </div>
        <div className="relative isolate rounded-xl sm:col-span-4 sm:h-full">
          <img
            aria-hidden
            alt="Starknet blockchain visualization"
            src="/golang_starknet_repo_banner.png"
            height={1000}
            width={1000}
            className="absolute inset-0 -z-10 rounded-2xl blur-xl"
          />
          <img
            alt="Starknet blockchain visualization"
            src="/golang_starknet_repo_banner.png"
            height={1000}
            width={1000}
            className="relative z-10 rounded-2xl"
          />
        </div>
      </div>
    </section>
  )
}

export default CallToAction
