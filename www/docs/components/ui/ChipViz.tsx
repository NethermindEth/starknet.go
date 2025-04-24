"use client"

import { motion } from "motion/react"
import { SolarMark } from "../../public/SolarMark"

const ChipViz = () => {
  const createVariants = ({
    scale,
    delay,
  }: {
    scale: number
    delay: number
  }) => ({
    initial: { scale: 1 },
    animate: {
      scale: [1, scale, 1],
      transition: {
        duration: 2,
        times: [0, 0.2, 1],
        ease: [0.23, 1, 0.32, 1],
        repeat: Infinity,
        repeatDelay: 2,
        delay,
      },
    },
  })

  return (
    <div className="relative flex items-center">
      <div className="relative">
        <motion.div
          variants={createVariants({ scale: 1.1, delay: 0 })}
          initial="initial"
          animate="animate"
          className="absolute -inset-px z-0 rounded-full bg-linear-to-r from-yellow-500 via-amber-500 to-orange-500 opacity-30 blur-xl"
        />
        <motion.div
          variants={createVariants({ scale: 1.08, delay: 0.1 })}
          initial="initial"
          animate="animate"
          className="relative z-0 min-h-[80px] min-w-[80px] rounded-full border bg-linear-to-b from-white to-orange-50 shadow-xl shadow-orange-500/20"
        >
          <motion.div
            variants={createVariants({ scale: 1.06, delay: 0.2 })}
            initial="initial"
            animate="animate"
            className="absolute inset-1 rounded-full bg-linear-to-t from-yellow-500 via-amber-500 to-orange-500 p-0.5 shadow-xl"
          >
            <div className="relative flex h-full w-full items-center justify-center overflow-hidden rounded-full bg-black/40 shadow-xs shadow-white/40 will-change-transform">
              <div className="size-full bg-black/30" />
              <motion.div
                variants={createVariants({ scale: 1.04, delay: 0.3 })}
                initial="initial"
                animate="animate"
                className="absolute inset-0 rounded-full bg-linear-to-t from-yellow-500 via-amber-500 to-orange-500 opacity-50 shadow-[inset_0_0_16px_4px_rgba(0,0,0,1)]"
              />
              <motion.div
                variants={createVariants({ scale: 1.02, delay: 0.4 })}
                initial="initial"
                animate="animate"
                className="absolute inset-[6px] rounded-full bg-white/10 p-1 backdrop-blur-[1px]"
              >
                <div className="relative flex h-full w-full items-center justify-center rounded-full bg-linear-to-br from-white to-gray-300 shadow-lg shadow-black/40">
                  <SolarMark className="w-6" />
                </div>
              </motion.div>
            </div>
          </motion.div>
        </motion.div>
      </div>
    </div>
  )
}

export default ChipViz
