import Banner from "./ui/Banner"
import { CallToAction } from "./ui/CallToAction"
import FeatureDivider from "./ui/FeatureDivider"
import Features from "./ui/Features"
import Footer from "./ui/Footer"
import { Hero } from "./ui/Hero"
import { Map } from "./ui/Map/Map"
import { SolarAnalytics } from "./ui/SolarAnalytics"
import Testimonial from "./ui/Testimonial"

export default function Home() {
  return (
    <main className="relative mx-auto flex flex-col w-full">
        <Banner/>
        <Hero />

      {/* <div className="mt-52 px-4 xl:px-0">
        <Features />
      </div>
      {/* <div className="mt-32 px-4 xl:px-0">
        <Testimonial />
      </div> */}
      {/* <FeatureDivider className="my-16 max-w-6xl" /> */}
      {/* <div className="px-4 xl:px-0">
        <Map />
      </div> */}
      {/* <FeatureDivider className="my-16 max-w-6xl" /> */}
      {/* <div className="mt-12 mb-40 px-4 xl:px-0">
        <SolarAnalytics />
      </div> */}
      {/* <div className="mt-10 mb-40 px-4 xl:px-0">
        <CallToAction />
      </div> */}
      {/* <div className="mt-10 mb-40 px-4 xl:px-0">
        <Footer />
      </div>  */}
    </main>
  )
}
