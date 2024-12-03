'use client'
import React, { useEffect, useState } from 'react'
import { fetchUserMoney } from '@/service/common-extend'
import type { UserMoney } from '@/models/common-extend'

const AccountMoneyExtend = () => {
  const [userMoney, setUserMoney] = useState<UserMoney>({ used_quota: 0, total_quota: 15 }) // TODO total_quota初始总额度
  const [isFetched, setIsFetched] = useState(false)
  const getUserMoney = async () => {
    const data: any = await fetchUserMoney()
    setUserMoney(data)
  }
  useEffect(() => {
    getUserMoney()
    setIsFetched(true)
  }, [])

  if (!isFetched)
    return null

  return (
    <div
      rel='noopener noreferrer'
      className='flex items-center leading-[18px] border border-gray-200 rounded-md text-xs text-gray-700 font-semibold overflow-hidden'>
      <div className='flex items-center px-2 py-1 bg-gray-100'>
         额度
      </div>
      <div className='px-2 py-1 bg-white border-l border-gray-200'>$ {`${userMoney.used_quota}`} / $ {`${userMoney.total_quota}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}</div>
    </div>
  )
}

export default AccountMoneyExtend
