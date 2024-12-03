'use client'
import {
  useEffect,
  useState,
} from 'react'
import { useTranslation } from 'react-i18next'
import { PlusIcon, XMarkIcon } from '@heroicons/react/20/solid'
import useSWR, { useSWRConfig } from 'swr'
import copy from 'copy-to-clipboard'
import SecretKeyGenerateModal from './secret-key-generate'
import SecretKeyQuotaSetExtendModal from './secret-key-quota-set-modal-extend' // 二开部分 - 密钥额度
import s from './style.module.css'
import Modal from '@/app/components/base/modal'
import Button from '@/app/components/base/button'
import {
  createApikey as createAppApikey,
  delApikey as delAppApikey,
  editApikey,
  fetchApiKeysList as fetchAppApiKeysList,
} from '@/service/apps'
import {
  createApikey as createDatasetApikey,
  delApikey as delDatasetApikey,
  fetchApiKeysList as fetchDatasetApiKeysList,
} from '@/service/datasets'
import type { ApiKeyItemResponse, ApikeyItemResponseWithQuotaLimitExtend, CreateApiKeyResponse } from '@/models/app' // 二开部分End - 密钥额度
import Tooltip from '@/app/components/base/tooltip'
import Loading from '@/app/components/base/loading'
import Confirm from '@/app/components/base/confirm'
import useTimestamp from '@/hooks/use-timestamp'
import { useAppContext } from '@/context/app-context'

type ISecretKeyModalProps = {
  isShow: boolean
  appId?: string
  onClose: () => void
}

const SecretKeyModal = ({
  isShow = false,
  appId,
  onClose,
}: ISecretKeyModalProps) => {
  const { t } = useTranslation()
  const { formatTime } = useTimestamp()
  const { currentWorkspace, isCurrentWorkspaceManager, isCurrentWorkspaceEditor } = useAppContext()
  const [showConfirmDelete, setShowConfirmDelete] = useState(false)
  const [isVisible, setVisible] = useState(false)
  const [newKey, setNewKey] = useState<CreateApiKeyResponse | undefined>(undefined)
  // ---------------------- 二开部分Begin - 密钥额度 ----------------------
  const [isVisibleExtend, setVisibleExtend] = useState(false)
  const [keyItem, setKeyItem] = useState<ApikeyItemResponseWithQuotaLimitExtend>({
    created_at: '',
    id: '',
    last_used_at: '',
    token: '',
    description: '',
    day_limit_quota: -1,
    month_limit_quota: -1,
  })
  // 打开新增密钥额度编辑框
  const openSecretKeyQuotaSetModalExtend = async () => {
    setVisibleExtend(true)
    setKeyItem({
      created_at: '',
      id: '',
      last_used_at: '',
      token: '',
      description: '',
      day_limit_quota: -1,
      month_limit_quota: -1,
    })
  }
  // 打开编辑密钥额度编辑框
  const openSecretKeyQuotaEditModalExtend = async (api: ApiKeyItemResponse) => {
    setVisibleExtend(true)
    setKeyItem({
      created_at: api.created_at,
      id: api.id,
      last_used_at: api.last_used_at,
      token: api.token,
      description: api.description,
      day_limit_quota: api.day_limit_quota,
      month_limit_quota: api.month_limit_quota,
    })
  }
  // 设置密钥额度数据
  const handleSetKeyDataSetQuotas = (newKeyItems: ApikeyItemResponseWithQuotaLimitExtend) => {
    setKeyItem(newKeyItems)
  }
  // ---------------------- 二开部分End - 密钥额度 ----------------------
  const { mutate } = useSWRConfig()
  const commonParams = appId
    ? { url: `/apps/${appId}/api-keys`, params: {} }
    : { url: '/datasets/api-keys', params: {} }
  const fetchApiKeysList = appId ? fetchAppApiKeysList : fetchDatasetApiKeysList
  const { data: apiKeysList } = useSWR(commonParams, fetchApiKeysList)

  const [delKeyID, setDelKeyId] = useState('')

  const [copyValue, setCopyValue] = useState('')

  useEffect(() => {
    if (copyValue) {
      const timeout = setTimeout(() => {
        setCopyValue('')
      }, 1000)

      return () => {
        clearTimeout(timeout)
      }
    }
  }, [copyValue])

  const onDel = async () => {
    setShowConfirmDelete(false)
    if (!delKeyID)
      return

    const delApikey = appId ? delAppApikey : delDatasetApikey
    const params = appId
      ? { url: `/apps/${appId}/api-keys/${delKeyID}`, params: {} }
      : { url: `/datasets/api-keys/${delKeyID}`, params: {} }
    await delApikey(params)
    mutate(commonParams)
  }

  const onCreate = async () => {
    // 二开部分 - 密钥额度，新增body传值
    const params = appId
      ? {
        url: `/apps/${appId}/api-keys`,
        body: {
          description: keyItem.description,
          day_limit_quota: keyItem.day_limit_quota,
          month_limit_quota: keyItem.month_limit_quota,
        },
      }
      : { url: '/datasets/api-keys', body: {} }
    const createApikey = appId ? createAppApikey : createDatasetApikey
    const res = await createApikey(params)
    setVisible(true)
    setVisibleExtend(false) // 二开部分 - 密钥额度，关闭创建额度编辑框
    setNewKey(res)
    mutate(commonParams)
  }

  // 二开部分 Begin - 密钥额度限制编辑
  const onEdit = async () => {
    const params = {
      url: `/apps/${appId}/api-keys`,
      body: {
        id: keyItem.id,
        description: keyItem.description,
        day_limit_quota: keyItem.day_limit_quota,
        month_limit_quota: keyItem.month_limit_quota,
      },
    }
    const res = await editApikey(params)
    setVisibleExtend(false)
    setNewKey(res)
    mutate(commonParams)
  }
  // 二开部分 Begin - 密钥额度限制编辑

  const generateToken = (token: string) => {
    return `${token.slice(0, 3)}...${token.slice(-20)}`
  }

  return (
    <Modal isShow={isShow} onClose={onClose} title={`${t('appApi.apiKeyModal.apiSecretKey')}`} className={`${s.customModalExtend} px-8 flex flex-col`}> {/* 二开部分 - 密钥额度限制，由customModal 改为 customModalExtend */}
      <XMarkIcon className={`w-6 h-6 absolute cursor-pointer text-gray-500 ${s.close}`} onClick={onClose} />
      <p className='mt-1 text-[13px] text-gray-500 font-normal leading-5 flex-shrink-0'>{t('appApi.apiKeyModal.apiSecretKeyTips')}</p>
      {!apiKeysList && <div className='mt-4'><Loading /></div>}
      {
        !!apiKeysList?.data?.length && (
          <div className='flex flex-col flex-grow mt-4 overflow-hidden'>
            <div className='flex items-center flex-shrink-0 text-xs font-semibold text-gray-500 border-b border-solid h-9'>
              <div className='flex-shrink-0 w-64 px-3'>{t('appApi.apiKeyModal.secretKey')}</div>
              <div className='flex-shrink-0 px-3 w-[150px]'>{t('appApi.apiKeyModal.created')}</div>
              <div className='flex-shrink-0 px-3 w-[150px]'>{t('appApi.apiKeyModal.lastUsed')}</div>
              {/* ---------------------- 二开部分Begin - 密钥额度限制 ---------------------- */}
              <div className='flex-shrink-0 px-3 w-[100px]'>{t('extend.apiKeyModal.descriptionPlaceholder')}</div>
              <div className='flex-shrink-0 px-3 w-[200px]'>{t('extend.apiKeyModal.dayLimit')}</div>
              <div className='flex-shrink-0 px-3 w-[200px]'>{t('extend.apiKeyModal.monthLimit')}</div>
              <div className='flex-shrink-0 px-3 w-[200px]'>{t('extend.apiKeyModal.accumulatedLimit')}</div>
              {/* ---------------------- 二开部分End - 密钥额度限制 ---------------------- */}
              <div className='flex-grow px-3'></div>
            </div>
            <div className='flex-grow overflow-auto'>
              {apiKeysList.data.map(api => (
                <div className='flex items-center text-sm font-normal text-gray-700 border-b border-solid h-9' key={api.id}>
                  <div className='flex-shrink-0 w-64 px-3 font-mono truncate'>{generateToken(api.token)}</div>
                  <div className='flex-shrink-0 px-3 truncate w-[150px]'>{formatTime(Number(api.created_at), t('appLog.dateTimeFormat') as string)}</div>{/* 二开部分 - 密钥额度限制，调整宽度 200 改为 150 ---------------------- */}
                  <div className='flex-shrink-0 px-3 truncate w-[150px]'>{api.last_used_at ? formatTime(Number(api.created_at), t('appLog.dateTimeFormat') as string) : t('appApi.never')}</div>{/* 二开部分 - 密钥额度限制，调整宽度 200 改为 150 ---------------------- */}
                  {/* ---------------------- 二开部分Begin - 密钥额度限制 ---------------------- */}
                  <div
                    className='flex-shrink-0 px-3 truncate w-[100px]'>{api.description}</div>
                  <div
                    className='flex-shrink-0 px-3 truncate w-[200px]'>$ {api.day_used_quota} / {api.day_limit_quota === -1 ? t('extend.apiKeyModal.noLimit') : `$ ${api.day_limit_quota}`}</div>
                  <div
                    className='flex-shrink-0 px-3 truncate w-[200px]'>$ {api.month_used_quota} / {api.month_limit_quota === -1 ? t('extend.apiKeyModal.noLimit') : `$ ${api.month_limit_quota}`}</div>
                  <div
                    className='flex-shrink-0 px-3 truncate w-[200px]'>$ {api.accumulated_quota}</div>
                  {/* ---------------------- 二开部分End - 密钥额度限制 ---------------------- */}
                  <div className='flex flex-grow px-3'>
                    <Tooltip
                      popupContent={copyValue === api.token ? `${t('appApi.copied')}` : `${t('appApi.copy')}`}
                      popupClassName='mr-1'
                    >
                      <div className={`flex items-center justify-center flex-shrink-0 w-6 h-6 mr-1 rounded-lg cursor-pointer hover:bg-gray-100 ${s.copyIcon} ${copyValue === api.token ? s.copied : ''}`} onClick={() => {
                        // setIsCopied(true)
                        copy(api.token)
                        setCopyValue(api.token)
                      }}></div>
                    </Tooltip>
                    {isCurrentWorkspaceManager
                      && <div className={`flex items-center justify-center flex-shrink-0 w-6 h-6 rounded-lg cursor-pointer ${s.trashIcon}`} onClick={() => {
                        setDelKeyId(api.id)
                        setShowConfirmDelete(true)
                      }}>
                      </div>
                    }
                    {/* // 二开部分 End - 密钥额度限制编辑 */}
                    {isCurrentWorkspaceManager
                      && <div className={`flex items-center justify-center flex-shrink-0 w-6 h-6 rounded-lg cursor-pointer ${s.editIcon}`} onClick={() => {
                        openSecretKeyQuotaEditModalExtend(api)
                      }}>
                      </div>
                    }
                    {/* // 二开部分 Begin - 密钥额度限制编辑 */}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )
      }
      <div className='flex'>
        <Button className={`flex flex-shrink-0 mt-4 ${s.autoWidth}`} onClick={openSecretKeyQuotaSetModalExtend} disabled={ !currentWorkspace || !isCurrentWorkspaceManager}> {/* 二开部分Begin - 密钥额度限制由onClick改为openSecretKeyQuotaSetModalExtend */}
          <PlusIcon className='flex flex-shrink-0 w-4 h-4' />
          <div className='text-xs font-medium text-gray-800'>{t('appApi.apiKeyModal.createNewSecretKey')}</div>
        </Button>
      </div>
      <SecretKeyGenerateModal className='flex-shrink-0' isShow={isVisible} onClose={() => setVisible(false)} newKey={newKey} />
      {showConfirmDelete && (
        <Confirm
          title={`${t('appApi.actionMsg.deleteConfirmTitle')}`}
          content={`${t('appApi.actionMsg.deleteConfirmTips')}`}
          isShow={showConfirmDelete}
          onConfirm={onDel}
          onCancel={() => {
            setDelKeyId('')
            setShowConfirmDelete(false)
          }}
        />
      )}

      {/* ----------------------二开部分Begin - 密钥额度限制---------------------- */}
      <SecretKeyQuotaSetExtendModal className='flex-shrink-0' isShow={isVisibleExtend} onClose={() => setVisibleExtend(false)} newKey={keyItem} onChange={handleSetKeyDataSetQuotas} onCreate={keyItem.id==''? onCreate:onEdit}/>
      {/* ----------------------二开部分End - 密钥额度限制---------------------- */}
    </Modal >
  )
}

export default SecretKeyModal
