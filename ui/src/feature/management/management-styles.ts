/** Tailwind classes matching Bootstrap management.html / management.css */
export const managementShellClass: string =
  'mx-auto my-8 w-full max-w-[1140px] px-3 2xl:max-w-[1320px]'

export const managementCardsColumnClass: string =
  'col-span-1 flex min-w-0 flex-col gap-3'

export const managementPageTitleClass: string =
  'text-[calc(1.375rem+1.5vw)] font-medium leading-tight text-[#212529]'

export const managementCardClass: string =
  'w-full gap-0 rounded-[6px] border border-[#dee2e6] bg-white py-0 text-base text-[#212529] shadow-none ring-0'

export const managementCardHeaderClass: string =
  'rounded-t-[5px] border-b border-[#dee2e6] bg-[rgba(33,37,41,0.03)] !pt-[8px] !pb-[8px] !pl-[16px] !pr-[16px]'

export const managementCardTitleClass: string =
  'text-base font-medium text-[#212529]'

export const managementCardContentClass: string = '!p-[16px]'

export const managementFieldTitleClass: string =
  'text-base font-medium text-[#212529]'

export const managementFieldTextClass: string =
  'mt-2.5 text-base text-[#212529]'

export const managementBootstrapButtonClass: string =
  'h-auto min-h-[38px] rounded-[0.375rem] px-3 py-1.5 text-base font-normal leading-normal'

/** Bootstrap `btn-sm` — padding 0.25rem 0.5rem, font-size 0.875rem (~31px tall).
 *  Includes overrides so it wins over `variant="bootstrap"` / default size classes. */
export const managementBootstrapSmButtonClass: string =
  'h-auto min-h-0 rounded-[0.25rem] px-2 py-1 text-sm font-normal leading-[1.5] gap-0'

export const managementBootstrapInputClass: string =
  'h-auto min-h-[38px] rounded-[0.375rem] border !border-[#dee2e6] !bg-white px-3 py-1.5 text-base leading-normal !text-[#212529] shadow-none md:text-base'

export const managementBootstrapLabelClass: string =
  'text-base font-normal text-[#212529]'

export const managementTabTriggerClass: string =
  'flex-none rounded-none !px-4 !py-2 text-base text-[#0d6efd] data-active:text-[#212529]'

export const managementCollateralButtonsClass: string = 'mt-[15px]'

export const managementLoadingBarClass: string = 'management-loading-bar'

/** Bootstrap danger `#dc3545` — invalid-feedback / table error text */
export const managementDangerTextClass: string = 'text-[#dc3545]'

/** Inline field error (invalid-feedback parity) — full static string for Tailwind JIT */
export const managementFieldErrorClass: string = 'mt-1 text-sm text-[#dc3545]'

/** Bootstrap danger button fill + hover `#bb2d3b` */
export const managementDangerButtonClass: string =
  'bg-[#dc3545] text-white hover:bg-[#bb2d3b]'

/** Bootstrap muted / form-text — `#6c757d` for secondary buttons; form-text uses secondary-color */
export const managementMutedTextClass: string = 'text-[#6c757d]'

/** Bootstrap `.form-text` — 0.875rem, secondary-color, margin-top 0.25rem */
export const managementFormTextClass: string =
  'mt-1 text-sm font-normal leading-[1.5] text-[rgba(33,37,41,0.75)]'

/** Bootstrap `.modal-title` / `h5` — 1.25rem, weight 500, line-height 1.5 */
export const managementModalTitleClass: string =
  'text-[1.25rem] leading-[1.5] font-medium text-[#212529]'

/** Bootstrap `.modal-header` */
export const managementModalHeaderClass: string =
  'flex-row items-center justify-between gap-0 space-y-0 border-b border-[#dee2e6] px-4 py-4'

/** Bootstrap `.modal-body` field group `.mb-3` */
export const managementModalFieldGroupClass: string = 'mb-3 grid gap-0'

/** Bootstrap secondary button `#6c757d` / hover `#5c636a` */
export const managementSecondaryButtonClass: string =
  'border-transparent bg-[#6c757d] text-white hover:bg-[#5c636a]'
