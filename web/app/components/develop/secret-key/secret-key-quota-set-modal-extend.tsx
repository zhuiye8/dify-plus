'use client'
import { useTranslation } from 'react-i18next'
import { XMarkIcon } from '@heroicons/react/20/solid'
import s from './style.module.css'
import Button from '@/app/components/base/button'
import Modal from '@/app/components/base/modal'
import type { ApikeyItemResponseWithQuotaLimitExtend } from '@/models/app'
import DayLimitItemExtend from '@/app/components/base/param-item/day-limit-item-extend'
import MonthLimitItemExtend from '@/app/components/base/param-item/month-limit-item-extend'

type ISecretKeyGenerateModalProps = {
  isShow: boolean
  onClose: () => void
  onCreate: () => void
  onChange: (keyItem: ApikeyItemResponseWithQuotaLimitExtend) => void
  newKey: ApikeyItemResponseWithQuotaLimitExtend
  className?: string
}

const SecretKeyQuotaSetExtendModal = ({
  isShow = false,
  onClose,
  onCreate,
  onChange,
  newKey,
  className,
}: ISecretKeyGenerateModalProps) => {
  const { t } = useTranslation()

  const handleParamChange = (key: string, value: any) => {
    if (key === 'day_limit_quota') {
      onChange({
        ...newKey,
        day_limit_quota: value,
      })
    }
    else if (key === 'month_limit_quota') {
      onChange({
        ...newKey,
        month_limit_quota: value,
      })
    }
    else if (key === 'description') {
      onChange({
        ...newKey,
        description: value,
      })
    }
  }

  const handleParamChangeDesc = (value: string) => {
    handleParamChange('description', value)
  }

  return (
    <Modal isShow={isShow} onClose={onClose} title={(newKey?.id ? '编辑' : '创建')+ `${t('appApi.apiKeyModal.apiSecretKey')}`}
      className={`px-8 ${className}`}>
      <XMarkIcon className={`w-6 h-6 absolute cursor-pointer text-gray-500 ${s.close}`} onClick={onClose}/>
      <p className='mt-1 text-[13px] text-gray-500 font-normal leading-5'>{t('extend.apiKeyModal.apiSecretKeyTips')}</p>
      <div className='my-4'>
        <input
          value={newKey?.description ?? ''}
          onChange={e => handleParamChangeDesc(e.target.value)}
          placeholder={t('extend.apiKeyModal.descriptionPlaceholder') || '密钥用途'}
          className='grow h-10 px-3 text-sm font-normal bg-gray-100 rounded-lg border border-transparent outline-none appearance-none caret-primary-600 placeholder:text-gray-400 hover:bg-gray-50 hover:border hover:border-gray-300 focus:bg-gray-50 focus:border focus:border-gray-300 focus:shadow-xs'
        />
      </div>
      <div className='my-4'>
        <DayLimitItemExtend
          value={newKey?.day_limit_quota ?? -1}
          onChange={handleParamChange}
          enable={true}
        />
      </div>
      <div className='my-4'>
        <MonthLimitItemExtend
          value={newKey?.month_limit_quota ?? -1}
          onChange={handleParamChange}
          enable={true}
        />
      </div>
      <div className='flex justify-end my-4'>
        <Button variant='primary' className={`flex-shrink-0 ${s.w64}`} onClick={onCreate}>
          {newKey?.id ? t('common.operation.save') : t('common.operation.create')}
        </Button>
      </div>

    </Modal>
  )
}

export default SecretKeyQuotaSetExtendModal
