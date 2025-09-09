import { describe, it, expect } from 'vitest'

import { mount } from '@vue/test-utils'
import IntroductionSection from '../IntroductionSection.vue'

describe('HomePage', () => {
  it('renders properly', () => {
    const wrapper = mount(IntroductionSection)
    expect(wrapper.text()).toContain(`Presenting Nunc's secret sauce for becoming a unicorn.`)
  })
})
