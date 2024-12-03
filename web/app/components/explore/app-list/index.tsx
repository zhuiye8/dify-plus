'use client'

import React, { useEffect, useState } from 'react' // 二开部份，新增useEffect
import { useRouter } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import { useContext } from 'use-context-selector'
import useSWR from 'swr'
import { useDebounceFn } from 'ahooks'
import Toast from '../../base/toast'
import s from './style.module.css'
import TagFilter from '@/app/components/base/tag-management/filter' // 二开部份
import cn from '@/utils/classnames'
import ExploreContext from '@/context/explore-context'
import type { App } from '@/models/explore'
import Category from '@/app/components/explore/category'
import AppCard from '@/app/components/explore/app-card'
import { fetchAppDetail, fetchAppList } from '@/service/explore'
import { importDSL } from '@/service/apps'
import { useTabSearchParams } from '@/hooks/use-tab-searchparams'
import CreateAppModal from '@/app/components/explore/create-app-modal'
import AppTypeSelector from '@/app/components/app/type-selector'
import type { CreateAppModalProps } from '@/app/components/explore/create-app-modal'
import Loading from '@/app/components/base/loading'
import { NEED_REFRESH_APP_LIST_KEY } from '@/config'
import { useAppContext } from '@/context/app-context'
import { getRedirection } from '@/utils/app-redirection'
// import Input from '@/app/components/base/input' // TODO 这里合并0.12.1版本时候多出来的未使用的，要看看二开部分对这个影响
import { DSLImportMode } from '@/models/app'
import SearchInput from '@/app/components/base/search-input' // Extend: Explore Add Search

type AppsProps = {
  pageType?: PageType
  onSuccess?: () => void
}

export enum PageType {
  EXPLORE = 'explore',
  CREATE = 'create',
}

const Apps = ({
  pageType = PageType.EXPLORE,
  onSuccess,
}: AppsProps) => {
  const { t } = useTranslation()
  const { isCurrentWorkspaceEditor } = useAppContext()
  const { push } = useRouter()
  const { hasEditPermission } = useContext(ExploreContext)
  const allCategoriesEn = t('explore.apps.allCategories', { lng: 'en' })

  // extend: start
  //
  // const [keywords, setKeywords] = useState('')
  // const [searchKeywords, setSearchKeywords] = useState('')
  //
  // const { run: handleSearch } = useDebounceFn(() => {
  //   setSearchKeywords(keywords)
  // }, { wait: 500 })
  //
  // const handleKeywordsChange = (value: string) => {
  //   setKeywords(value)
  //   handleSearch()
  // }
  // extend: stop

  const [currentType, setCurrentType] = useState<string>('')
  const [currCategory, setCurrCategory] = useTabSearchParams({
    defaultTab: allCategoriesEn,
    disableSearchParams: pageType !== PageType.EXPLORE,
  })

  const {
    data: { categories, allList },
  } = useSWR(
    ['/explore/apps'],
    () =>
      fetchAppList().then(({ categories, recommended_apps }) => ({
        categories,
        allList: recommended_apps.sort((a, b) => a.position - b.position),
      })),
    {
      fallbackData: {
        categories: [],
        allList: [],
      },
    },
  )

  // -------------- start: Filter List Extension ---------------

  const [filteredListExtend, setFilteredListExtend] = useState<any[]>([])
  const [tagFilterValue, setTagFilterValue] = useState<string[]>([])
  const [keywordsValue, setKeywordsValue] = useState<string>('')

  // Explore Add Search

  useEffect(() => {
    const newList = []
    let cacheList = []
    const idList: string[] = []
    if (currCategory === allCategoriesEn) {
      if (!currentType)
        cacheList = allList
      else if (currentType === 'chatbot')
        cacheList = allList.filter(item => (item.app.mode === 'chat' || item.app.mode === 'advanced-chat'))
      else if (currentType === 'agent')
        cacheList = allList.filter(item => (item.app.mode === 'agent-chat'))
      else
        cacheList = allList.filter(item => (item.app.mode === 'workflow'))
    }
    else {
      if (!currentType)
        cacheList = allList.filter(item => item.category === currCategory)
      else if (currentType === 'chatbot')
        cacheList = allList.filter(item => (item.app.mode === 'chat' || item.app.mode === 'advanced-chat') && item.category === currCategory)
      else if (currentType === 'agent')
        cacheList = allList.filter(item => (item.app.mode === 'agent-chat') && item.category === currCategory)
      else
        cacheList = allList.filter(item => (item.app.mode === 'workflow') && item.category === currCategory)
    }
    // 循环遍历cacheList，去重
    for (const i in cacheList) {
      if (!idList.includes(cacheList[i].app_id)) {
        idList.push(cacheList[i].app_id)
        newList.push(cacheList[i])
      }
    }
    // 返回去重后的数组
    if (allList.length > 0)
      setFilteredListExtend(newList)
  }, [currentType, currCategory, allCategoriesEn, allList])

  const { run: handleSearch } = useDebounceFn(() => {
    const cacheList: any[] = []
    const idList: string[] = []
    for (const i in allList) {
      if (keywordsValue.length > 0) {
        if (!(allList[i].description.includes(keywordsValue) || allList[i].app.name.includes(keywordsValue)))
          continue
      }
      if (tagFilterValue.length > 0) {
        if (!tagFilterValue.includes(allList[i].category))
          continue
      }
      if (!idList.includes(allList[i].app_id)) {
        idList.push(allList[i].app_id)
        cacheList.push(allList[i])
      }
    }
    // save
    setFilteredListExtend(cacheList)
  }, { wait: 500 })
  const handleKeywordsChange = (value: string) => {
    setKeywordsValue(value)
    handleSearch()
  }
  const handleTagsChange = (value: string[]) => {
    setTagFilterValue(value)
    handleSearch()
  }
  // -------------- stop: Filter List Extension ---------------

  // const filteredList = useMemo(() => {
  //   if (currCategory === allCategoriesEn) {
  //     if (!currentType)
  //       return allList
  //     else if (currentType === 'chatbot')
  //       return allList.filter(item => (item.app.mode === 'chat' || item.app.mode === 'advanced-chat'))
  //     else if (currentType === 'agent')
  //       return allList.filter(item => (item.app.mode === 'agent-chat'))
  //     else
  //       return allList.filter(item => (item.app.mode === 'workflow'))
  //   }
  //   else {
  //     if (!currentType)
  //       return allList.filter(item => item.category === currCategory)
  //     else if (currentType === 'chatbot')
  //       return allList.filter(item => (item.app.mode === 'chat' || item.app.mode === 'advanced-chat') && item.category === currCategory)
  //     else if (currentType === 'agent')
  //       return allList.filter(item => (item.app.mode === 'agent-chat') && item.category === currCategory)
  //     else
  //       return allList.filter(item => (item.app.mode === 'workflow') && item.category === currCategory)
  //   }
  // }, [currentType, currCategory, allCategoriesEn, allList])
  //
  // const searchFilteredList = useMemo(() => {
  //   if (!searchKeywords || !filteredList || filteredList.length === 0)
  //     return filteredList
  //
  //   const lowerCaseSearchKeywords = searchKeywords.toLowerCase()
  //
  //   return filteredList.filter(item =>
  //     item.app && item.app.name && item.app.name.toLowerCase().includes(lowerCaseSearchKeywords),
  //   )
  // }, [searchKeywords, filteredList])

  const [currApp, setCurrApp] = React.useState<App | null>(null)
  const [isShowCreateModal, setIsShowCreateModal] = React.useState(false)
  const onCreate: CreateAppModalProps['onConfirm'] = async ({
    name,
    icon_type,
    icon,
    icon_background,
    description,
  }) => {
    const { export_data } = await fetchAppDetail(
      currApp?.app.id as string,
    )
    try {
      const app = await importDSL({
        mode: DSLImportMode.YAML_CONTENT,
        yaml_content: export_data,
        name,
        icon_type,
        icon,
        icon_background,
        description,
      })
      setIsShowCreateModal(false)
      Toast.notify({
        type: 'success',
        message: t('app.newApp.appCreated'),
      })
      if (onSuccess)
        onSuccess()
      localStorage.setItem(NEED_REFRESH_APP_LIST_KEY, '1')
      getRedirection(isCurrentWorkspaceEditor, { id: app.app_id }, push)
    }
    catch (e) {
      Toast.notify({ type: 'error', message: t('app.newApp.appCreateFailed') })
    }
  }

  if (!categories || categories.length === 0) {
    return (
      <div className="flex h-full items-center">
        <Loading type="area" />
      </div>
    )
  }

  return (
    <div className={cn(
      'flex flex-col',
      pageType === PageType.EXPLORE ? 'h-full border-l border-gray-200' : 'h-[calc(100%-56px)]',
    )}>
      {pageType === PageType.EXPLORE && (
        <div className='shrink-0 pt-6 px-12'>
          <div className={`mb-1 ${s.textGradient} text-xl font-semibold`}>{t('explore.apps.title')}</div>
          <div className='text-gray-500 text-sm'>{t('explore.apps.description')}</div>
        </div>
      )}
      <div className={cn(
        'flex items-center justify-between mt-6',
        pageType === PageType.EXPLORE ? 'px-12' : 'px-8',
      )}>
        <>
          {pageType !== PageType.EXPLORE && (
            <>
              <AppTypeSelector value={currentType} onChange={setCurrentType}/>
              <div className='mx-2 w-[1px] h-3.5 bg-gray-200'/>
            </>
          )}
          <Category
            list={categories}
            value={currCategory}
            onChange={setCurrCategory}
            allCategoriesEn={allCategoriesEn}
          />
        </>

        {/* extend: Application Center Search Start */}
        <div className={cn('flex items-center gap-2', s.rightSearch)}>
          <TagFilter type="app" value={tagFilterValue} onChange={handleTagsChange} defaultValue={categories} />
          <SearchInput className="w-[200px]" value={keywordsValue} onChange={handleKeywordsChange}/>
        </div>
        {/* extend: Application Center Search Stop */}

      </div>

      <div className={cn(
        'relative flex flex-1 pb-6 flex-col overflow-auto bg-gray-100 shrink-0 grow',
        pageType === PageType.EXPLORE ? 'mt-4' : 'mt-0 pt-2',
      )}>
        <nav
          className={cn(
            s.appList,
            'grid content-start shrink-0',
            pageType === PageType.EXPLORE ? 'gap-4 px-6 sm:px-12' : 'gap-3 px-8  sm:!grid-cols-2 md:!grid-cols-3 lg:!grid-cols-4',
          )}>
          {/* extend start */}
          {filteredListExtend.map(app => (
            <AppCard
              key={app.app_id}
              isExplore={pageType === PageType.EXPLORE}
              app={app}
              canCreate={hasEditPermission}
              onCreate={() => {
                setCurrApp(app)
                setIsShowCreateModal(true)
              }}
            />
          ))}
          {/* extend stop */}
        </nav>
      </div>
      {isShowCreateModal && (
        <CreateAppModal
          appIconType={currApp?.app.icon_type || 'emoji'}
          appIcon={currApp?.app.icon || ''}
          appIconBackground={currApp?.app.icon_background || ''}
          appIconUrl={currApp?.app.icon_url}
          appName={currApp?.app.name || ''}
          appDescription={currApp?.app.description || ''}
          show={isShowCreateModal}
          onConfirm={onCreate}
          onHide={() => setIsShowCreateModal(false)}
        />
      )}
    </div>
  )
}

export default React.memo(Apps)
