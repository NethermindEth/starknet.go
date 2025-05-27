import React, { useState, useEffect } from 'react';
import { AlertTriangle, Code, Zap, X } from 'lucide-react';

export default function DevelopmentBanner() {
  const [isVisible, setIsVisible] = useState(true);
  const [isAnimating, setIsAnimating] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => setIsAnimating(false), 3000);
    return () => clearTimeout(timer);
  }, []);

  if (!isVisible) return null;

  return (
    <div className="relative overflow-hidden bg-gradient-to-r from-amber-500 via-orange-500 to-red-500 text-white shadow-lg">
      {/* Animated background pattern */}
      <div className="absolute inset-0 opacity-20">
        <div className="absolute inset-0 bg-[repeating-linear-gradient(45deg,transparent,transparent_10px,rgba(255,255,255,0.1)_10px,rgba(255,255,255,0.1)_20px)]"></div>
      </div>
      
      {/* Animated dots */}
      <div className="absolute inset-0 overflow-hidden">
        {[...Array(6)].map((_, i) => (
          <div
            key={i}
            className={`absolute w-2 h-2 bg-white/30 rounded-full animate-bounce`}
            style={{
              left: `${15 + i * 15}%`,
              top: '50%',
              animationDelay: `${i * 0.2}s`,
              animationDuration: '2s'
            }}
          />
        ))}
      </div>

      <div className="relative flex items-center justify-center px-4 py-3 sm:px-6">
        <div className="flex items-center space-x-3">
          {/* Pulsing warning icon */}
          <div className="relative">
            <AlertTriangle className={`w-5 h-5 ${isAnimating ? 'animate-pulse' : ''}`} />
            <div className="absolute -top-1 -right-1 w-3 h-3 bg-white rounded-full animate-ping opacity-75"></div>
          </div>

          {/* Main text with typewriter effect */}
          <div className="flex items-center space-x-2 font-semibold text-sm sm:text-base">
            <span className="hidden sm:inline">⚠️</span>
            <span className={isAnimating ? 'animate-pulse' : ''}>
              Documentation Under Active Development
            </span>
          </div>

          {/* Animated icons */}
          <div className="flex items-center space-x-2">
            <Code className="w-4 h-4 animate-spin" style={{ animationDuration: '3s' }} />
            <Zap className={`w-4 h-4 ${isAnimating ? 'animate-bounce' : ''}`} />
          </div>
        </div>

        {/* Close button */}
        <button
          onClick={() => setIsVisible(false)}
          className="absolute right-2 sm:right-4 p-1 hover:bg-white/20 rounded-full transition-colors duration-200 group"
          aria-label="Dismiss banner"
        >
          <X className="w-4 h-4 group-hover:scale-110 transition-transform duration-200" />
        </button>
      </div>

      {/* Bottom accent line */}
      <div className="h-1 bg-gradient-to-r from-transparent via-white/50 to-transparent"></div>
      
      {/* Shimmer effect */}
      <div className="absolute inset-0 -skew-x-12 bg-gradient-to-r from-transparent via-white/20 to-transparent w-1/4 animate-pulse" 
           style={{ 
             animation: 'shimmer 3s ease-in-out infinite',
             transform: 'translateX(-100%) skewX(-12deg)'
           }}>
      </div>

      <style>{`
        @keyframes shimmer {
          0% { transform: translateX(-100%) skewX(-12deg); }
          100% { transform: translateX(400%) skewX(-12deg); }
        }
      `}</style>
    </div>
  );
}