import React, { useEffect, useRef, useState } from 'react'
import mermaid from 'mermaid'
import { usePrevious } from 'ahooks'
import CryptoJS from 'crypto-js'
import { ExclamationTriangleIcon } from '@heroicons/react/24/outline'
import type { Property } from 'csstype' // Markdown embedded images support click-to-zoom, providing users with a better user experience
import cn from 'classnames' // Markdown embedded images support click-to-zoom, providing users with a better user experience
import s from './modal.module.css' // Markdown embedded images support click-to-zoom, providing users with a better user experience
import LoadingAnim from '@/app/components/base/chat/chat/loading-anim'

let mermaidAPI: any
mermaidAPI = null

if (typeof window !== 'undefined') {
  mermaid.initialize({
    startOnLoad: true,
    theme: 'default',
    flowchart: {
      htmlLabels: true,
      useMaxWidth: true,
    },
  })
  mermaidAPI = mermaid.mermaidAPI
}

const style = {
  minWidth: '480px',
  height: 'auto',
  overflow: 'auto',
}

const svgToBase64 = (svgGraph: string) => {
  const svgBytes = new TextEncoder().encode(svgGraph)
  const blob = new Blob([svgBytes], { type: 'image/svg+xml;charset=utf-8' })
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onloadend = () => resolve(reader.result)
    reader.onerror = reject
    reader.readAsDataURL(blob)
  })
}

// eslint-disable-next-line react/display-name
const Flowchart = React.forwardRef((props: {
  PrimitiveCode: string
}, ref) => {
  const [svgCode, setSvgCode] = useState(null)
  const chartId = useRef(`flowchart_${CryptoJS.MD5(props.PrimitiveCode).toString()}`)
  const prevPrimitiveCode = usePrevious(props.PrimitiveCode)
  const [isLoading, setIsLoading] = useState(true)
  const timeRef = useRef<NodeJS.Timeout>()
  const [errMsg, setErrMsg] = useState('')
  // Extend: Start Markdown embedded images support click-to-zoom, providing users with a better user experience
  const [showModal, setShowModal] = useState(false)
  const [imgStyle, setImgStyle] = useState({ width: '0px', height: '0px' })
  const [imgBarStyle, setImgBarStyle] = useState({ width: '0px', height: '0px', overflowY: 'auto' as Property.OverflowY })
  // Extend: Stop Markdown embedded images support click-to-zoom, providing users with a better user experience

  const renderFlowchart = async (PrimitiveCode: string) => {
    try {
      if (typeof window !== 'undefined' && mermaidAPI) {
        const svgGraph = await mermaidAPI.render(chartId.current, PrimitiveCode)
        const base64Svg: any = await svgToBase64(svgGraph.svg)
        setSvgCode(base64Svg)
        setIsLoading(false)
        if (chartId.current && base64Svg)
          localStorage.setItem(chartId.current, base64Svg)
      }
    }
    catch (error) {
      if (prevPrimitiveCode === props.PrimitiveCode) {
        setIsLoading(false)
        setErrMsg((error as Error).message)
      }
    }
  }

  // Extend: Start Markdown embedded images support click-to-zoom, providing users with a better user experience
  const openImage = async () => {
    if (!window)
      return
    const img = new Image()
    let winWidth = window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth
    let winHeight = window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight
    img.onload = () => {
      setShowModal(true)
      winWidth = winWidth * 0.9
      winHeight = winHeight * 0.9
      winWidth = winWidth > 1920 ? 1920 : winWidth
      const multiple = winWidth / img.width
      const imgHeight = img.height * multiple
      if (imgHeight < winHeight)
        winHeight = imgHeight
      setImgBarStyle({ width: `${winWidth}px`, height: `${winHeight}px`, overflowY: 'auto' })
      setImgStyle({ width: `${img.width * multiple}px`, height: `${imgHeight}px` })
    }
    if (svgCode)
      img.src = svgCode
  }
  // clone
  const closeImage = () => {
    setShowModal(false)
  }
  // Extend: Stop Markdown embedded images support click-to-zoom, providing users with a better user experience

  useEffect(() => {
    const cachedSvg: any = localStorage.getItem(chartId.current)
    if (cachedSvg) {
      setSvgCode(cachedSvg)
      setIsLoading(false)
      return
    }
    if (timeRef.current)
      clearTimeout(timeRef.current)

    timeRef.current = setTimeout(() => {
      renderFlowchart(props.PrimitiveCode)
    }, 300)
  }, [props.PrimitiveCode])

  return (
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    <div ref={ref}>
      {
        svgCode
        && <div className="mermaid" style={style}>
          {svgCode && <img onClick={openImage} src={svgCode} style={{ width: '100%', height: 'auto' }} alt="Mermaid chart" />}
        </div>
      }
      {isLoading
        && <div className='py-4 px-[26px]'>
          <LoadingAnim type='text' />
        </div>
      }
      {
        errMsg
        && <div className='py-4 px-[26px]'>
          <ExclamationTriangleIcon className='w-6 h-6 text-red-500' />
          &nbsp;
          {errMsg}
        </div>
      }
      {/* Extend: Start Markdown embedded images support click-to-zoom, providing users with a better user experience */}
      {
        (showModal && imgStyle && svgCode) && <div onClick={closeImage}>
          <div className={cn(s.mask_body)}>
            <div className={cn(s.mask_layer)}></div>
            <div className={cn(s.mask_dialog)} style={imgBarStyle}>
              <img
                alt="result image preview"
                src={svgCode}
                style={imgStyle}
                className={cn('h-full', 'w-full', 'object-contain', s.mask_layer)}
              />
            </div>
          </div>
        </div>
      }
      {/* Extend: Stop Markdown embedded images support click-to-zoom, providing users with a better user experience */}
    </div>
  )
})

export default Flowchart
