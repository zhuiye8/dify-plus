'use client'
import { Menu, Transition } from '@headlessui/react'
import { RiArrowUpSLine } from '@remixicon/react'
import React, { useEffect, Fragment } from 'react'
import cn from 'classnames'

type ISelectProps = {
  items: Array<{ value: string; name: string }>
  value?: string
  className?: string
  menuClass?: string
  isDown?: boolean
  onChange?: (value: string) => void
}

export default function Select({
  items,
  value,
  isDown,
  menuClass,
  onChange,
}: ISelectProps) {
  const [open, setOpen] = React.useState(false)
  const [time, setTime] = React.useState(Number(new Date().getTime()))
  const item = items.find(item => item.value === value)
  
  useEffect(() => {
    setTime(Number(new Date().getTime()))
  }, [open])

  return (
    <div className={cn('w-56 text-right', menuClass)} onBlur={(e) => {
      e.stopPropagation()
      e.preventDefault()
      setTimeout(() => {
        if ((Number(new Date().getTime()) - time) > 2000) {
          setOpen(!open)
        }
      }, 800)
    }}>
      <Menu as="div" className="relative inline-block text-left">
        <div>
          <Menu.Button className="inline-flex w-full h-[44px]justify-center items-center
          rounded-lg px-[10px] py-[6px]
          text-gray-900 text-[13px] font-medium
          border border-gray-200
          hover:bg-gray-100"
          onClick={(e) => {
            e.stopPropagation()
            e.preventDefault()
            setTime(Number(new Date().getTime()))
            setOpen(!open)
          }}
          >
            {item?.name}
            <RiArrowUpSLine className='w-3.5 text-gray-700'/>
          </Menu.Button>
        </div>
        <Transition
          show={open}
          as={Fragment}
          enter="transition ease-out duration-100"
          enterFrom="transform opacity-0 scale-95"
          enterTo="transform opacity-100 scale-100"
          leave="transition ease-in duration-75"
          leaveFrom="transform opacity-100 scale-100"
          leaveTo="transform opacity-0 scale-95"
        >
          <Menu.Items
            className="absolute right-0 mt-2 w-[200px] origin-top-right divide-y divide-gray-100 rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none z-10"
            style={{ bottom: isDown ? '2.5rem' : 'auto' }}
          >
            <div className="px-1 py-1 ">
              {items.map((item, index) => {
                return <Menu.Item key={index}>
                  {({ active }) => (
                    <button
                      className={`${active ? 'bg-gray-100' : ''
                      } group flex w-full items-center rounded-lg px-3 py-2 text-sm text-gray-700`}
                      onClick={(evt) => {
                        evt.stopPropagation()
                        evt.preventDefault()
                        setTime(Number(new Date().getTime()))
                        setOpen(false)
                        onChange && onChange(item.value)
                      }}
                    >
                      {item.name}
                    </button>
                  )}
                </Menu.Item>
              })}

            </div>

          </Menu.Items>
        </Transition>
      </Menu>
    </div>
  )
}
