"use client"
import { useEffect, useRef } from "react"
import dynamic from 'next/dynamic'

type Particle = {
  x: number
  y: number
  vx: number
  vy: number
  size: number
  opacity: number
  maxOpacity: number
  life: number
  maxLife: number
}

const CleanBackground = () => {
  const canvasRef = useRef<HTMLCanvasElement>(null)

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext("2d")
    if (!ctx) return

    console.log("Canvas initialized, starting clean animation...")

    let animationFrameId: number
    let frameCount = 0
    
    const particles: Particle[] = []
    const maxParticles = 40 // More particles
    
    // Center exclusion zone
    const centerX = canvas.width / 2
    const centerY = canvas.height / 2
    const exclusionRadius = 220 // Slightly smaller exclusion zone

    const createParticle = (): Particle => {
      let x, y
      
      // Keep trying until we get a position outside the center area
      do {
        x = Math.random() * canvas.width
        y = Math.random() * canvas.height
      } while (
        Math.sqrt((x - centerX) ** 2 + (y - centerY) ** 2) < exclusionRadius
      )
      
      return {
        x,
        y,
        vx: (Math.random() - 0.5) * 0.8, // Slightly faster movement
        vy: (Math.random() - 0.5) * 0.8,
        size: Math.random() * 3 + 0.5, // Varied sizes
        opacity: 0,
        maxOpacity: Math.random() * 0.4 + 0.15, // Slightly more visible
        life: 0,
        maxLife: Math.random() * 400 + 300 // Longer life
      }
    }

    // Initialize particles
    for (let i = 0; i < maxParticles; i++) {
      particles.push(createParticle())
    }

    const draw = () => {
      frameCount++
      
      // Clear canvas
      ctx.clearRect(0, 0, canvas.width, canvas.height)
      
      // Update and draw particles
      for (let i = particles.length - 1; i >= 0; i--) {
        const p = particles[i]
        
        // Update position
        p.x += p.vx
        p.y += p.vy
        p.life++
        
        // Fade in and out
        if (p.life < p.maxLife * 0.2) {
          p.opacity = (p.life / (p.maxLife * 0.2)) * p.maxOpacity
        } else if (p.life > p.maxLife * 0.8) {
          p.opacity = (1 - (p.life - p.maxLife * 0.8) / (p.maxLife * 0.2)) * p.maxOpacity
        } else {
          p.opacity = p.maxOpacity
        }
        
        // Remove dead particles
        if (p.life >= p.maxLife) {
          particles.splice(i, 1)
          particles.push(createParticle())
          continue
        }
        
        // Wrap around edges
        if (p.x < 0) p.x = canvas.width
        if (p.x > canvas.width) p.x = 0
        if (p.y < 0) p.y = canvas.height
        if (p.y > canvas.height) p.y = 0
        
        // Don't draw if too close to center
        const distFromCenter = Math.sqrt((p.x - centerX) ** 2 + (p.y - centerY) ** 2)
        if (distFromCenter < exclusionRadius) {
          continue
        }
        
        // Draw particle with subtle glow
        const glowSize = p.size + 1
        
        // Outer glow
        ctx.fillStyle = `rgba(10, 8, 55, ${p.opacity * 0.3})`
        ctx.beginPath()
        ctx.arc(p.x, p.y, glowSize, 0, Math.PI * 2)
        ctx.fill()
        
        // Main particle
        ctx.fillStyle = `rgba(10, 8, 55, ${p.opacity})`
        ctx.beginPath()
        ctx.arc(p.x, p.y, p.size, 0, Math.PI * 2)
        ctx.fill()
      }
      
      // Draw subtle connecting lines between nearby particles
      for (let i = 0; i < particles.length; i++) {
        for (let j = i + 1; j < particles.length; j++) {
          const p1 = particles[i]
          const p2 = particles[j]
          
          const dx = p1.x - p2.x
          const dy = p1.y - p2.y
          const distance = Math.sqrt(dx * dx + dy * dy)
          
          // Only connect if both particles are visible and close
          if (distance < 100 && p1.opacity > 0.1 && p2.opacity > 0.1) {
            // Check if line would cross center exclusion zone
            const midX = (p1.x + p2.x) / 2
            const midY = (p1.y + p2.y) / 2
            const distFromCenter = Math.sqrt((midX - centerX) ** 2 + (midY - centerY) ** 2)
            
            if (distFromCenter > exclusionRadius) {
              const lineOpacity = (1 - distance / 100) * Math.min(p1.opacity, p2.opacity) * 0.2
              ctx.strokeStyle = `rgba(10, 8, 55, ${lineOpacity})`
              ctx.lineWidth = 0.5
              ctx.beginPath()
              ctx.moveTo(p1.x, p1.y)
              ctx.lineTo(p2.x, p2.y)
              ctx.stroke()
            }
          }
        }
      }
      
      animationFrameId = requestAnimationFrame(draw)
    }

    draw()

    return () => {
      cancelAnimationFrame(animationFrameId)
    }
  }, [])

  return (
    <canvas 
      ref={canvasRef} 
      width={1200} 
      height={800} 
      className="w-full h-full"
      style={{
        display: 'block',
        background: 'transparent'
      }}
    />
  )
}

// Export as dynamic component with no SSR
const GameOfLife = dynamic(() => Promise.resolve(CleanBackground), {
  ssr: false,
  loading: () => <div className="w-full h-full">Loading canvas...</div>
})

export default GameOfLife