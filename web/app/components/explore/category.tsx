'use client'
import type { FC } from 'react'
import React from 'react'
import { useTranslation } from 'react-i18next' // add cancelSyncToAppTemplate to list extend
import cn from '@/utils/classnames'
import { ThumbsUp } from '@/app/components/base/icons/src/vender/line/alertsAndFeedback'

export type ICategoryProps = {
  className?: string
  list: string[]
  value: string
  onChange: (value: string) => void
  /**
   * default value for search param 'category' in en
   */
  allCategoriesEn: string
}

const Category: FC<ICategoryProps> = ({
  className,
  list,
  value,
  onChange,
  allCategoriesEn,
}) => {
  const { t } = useTranslation()
  const isAllCategories = !list.includes(value as AppCategory) || value === allCategoriesEn

  const itemClassName = (isSelected: boolean) => cn(
    'flex items-center px-3 py-[7px] h-[32px] rounded-lg border-[0.5px] border-transparent text-gray-700 font-medium leading-[18px] cursor-pointer hover:bg-gray-200',
    isSelected && 'bg-white border-gray-200 shadow-xs text-primary-600 hover:bg-white',
  )

  return (
    <div className={cn(className, 'flex space-x-1 text-[13px] flex-wrap')}>
      <div
        className={itemClassName(isAllCategories)}
        onClick={() => onChange(allCategoriesEn)}
      >
        <ThumbsUp className='mr-1 w-3.5 h-3.5' />
        {t('explore.apps.allCategories')}
      </div>
      {list.filter(name => name !== allCategoriesEn).map(name => (
        <div
          key={name}
          className={itemClassName(name === value)}
          onClick={() => onChange(name)}
        >
          {name}
        </div>
      ))}
    </div>
  )
}

export default React.memo(Category)
