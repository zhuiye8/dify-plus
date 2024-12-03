'use client'
import React, { useState } from 'react'
import dayjs from 'dayjs'
import quarterOfYear from 'dayjs/plugin/quarterOfYear'
import { useTranslation } from 'react-i18next'
import {
  PeriodParams,
  WorkflowCostChart,
  WorkflowMessagesChart,
} from '@/app/components/app/overview/appChart'
import { AvgSessionInteractions, AvgUserInteractions, ConversationsChart, CostChart } from '@/app/components/app/overview/appChart'

import type { Item } from '@/app/components/base/select'
import { SimpleSelect } from '@/app/components/base/select'
import { TIME_PERIOD_MAPPING } from '@/app/components/app/log/filter'
import { useStore as useAppStore } from '@/app/components/app/store'

dayjs.extend(quarterOfYear)

const today = dayjs()

const queryDateFormat = 'YYYY-MM-DD HH:mm'

export type UserOverViewProps = {
  params: { appId: string }
}

const UserOverView = ({ params: { appId } }: UserOverViewProps) => {
  const { t } = useTranslation()
  const appDetail = useAppStore(state => state.appDetail)
  const model = appDetail?.mode
  const isChatApp = model !== 'completion' && model !== 'workflow'
  const [period, setPeriod] = useState<PeriodParams>({
    name: t('appLog.filter.period.last7days'),
    query: {
      start: today.subtract(7, 'day').startOf('day').format(queryDateFormat),
      end: today.format(queryDateFormat),
      account: true,
    },
  })

  const onSelect = (item: Item) => {
    if (item.value === 'all') {
      setPeriod({ name: item.name, query: undefined })
    }
    else if (item.value === 0) {
      const startOfToday = today.startOf('day').format(queryDateFormat)
      const endOfToday = today.endOf('day').format(queryDateFormat)
      setPeriod({ name: item.name, query: { start: startOfToday, end: endOfToday, account: true } })
    }
    else {
      setPeriod({
        name: item.name,
        query: {
          start: today.subtract(item.value as number, 'day').startOf('day').format(queryDateFormat),
          end: today.format(queryDateFormat),
          account: true,
        },
      })
    }
  }

  if (!appDetail)
    return null

  return (
    <div>
      <div className='flex flex-row items-center mt-8 mb-4 text-gray-900 text-base'>
        <span className='mr-3'>{t('appOverview.analysis.title')}</span>
        <SimpleSelect
          items={Object.entries(TIME_PERIOD_MAPPING).map(([k, v]) => ({ value: k, name: t(`appLog.filter.period.${v.name}`) }))}
          className='mt-0 !w-40'
          onSelect={onSelect}
          defaultValue={7}
        />
      </div>
      {model === 'workflow' && (
        <>
          {/* Extend: Workflow personal detection error */}
          <div className='grid gap-6 grid-cols-1 xl:grid-cols-2 w-full mb-6'>
            <WorkflowMessagesChart period={period} id={appId}/>
            <WorkflowCostChart period={period} id={appId}/>
          </div>
          <div className='grid gap-6 grid-cols-1 xl:grid-cols-2 w-full mb-6'>
            <AvgUserInteractions period={period} id={appId}/>
          </div>
        </>
      )}
      {model !== 'workflow' && (
        <>
          <div className='grid gap-6 grid-cols-1 xl:grid-cols-2 w-full mb-6'>
            <ConversationsChart period={period} id={appId}/>
            {model !== 'completion' && (isChatApp
              ? (
                <AvgSessionInteractions period={period} id={appId}/>
              )
              : (
                <AvgUserInteractions period={period} id={appId}/>
              ))}
          </div>
          <div className='grid gap-6 grid-cols-1 xl:grid-cols-2 w-full mb-6'>
            <CostChart period={period} id={appId}/>
          </div>
        </>
      )}
    </div>
  )
}

export default UserOverView
