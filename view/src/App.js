import React, { Component } from 'react'
import {debounce } from 'lodash'
import axios from 'axios'
import MarkdownInput from './MarkdownInput'
import SlatePreview from './OldApp'

class App extends Component {
  state = {
    md: "",
    slate: [],
    queue: [],
    pending: false
  }

  endpoint() {
    return "http://localhost:8080/md2slate"
  }

  setSlate = debounce((value) => {
    if (this.state.pending) {
      return this.setState({queue: [value]})
    }
    this.setState({pending: true})
    axios.post(this.endpoint(), {
      body: value
    }).then(resp => {
      this.setState({slate: resp.data, pending: false})
    }).catch(err => {
      this.setState({error: err.message, pending: false})
    }).then(() => {
      const {queue} = this.state
      if (queue.length > 0) {
        const next = queue[queue.length - 1]
        this.setState({queue: []}, () => {
          this.setSlate(next)
        })
      }
    })
  }, 200)

  setMD = value => {
    this.setState({md: value})
    this.setSlate(value)
  }

  render() {
    return (
      <div>
        <MarkdownInput value={this.state.md} onChange={this.setMD} />
        <SlatePreview value={this.state.slate} />
      </div>
    )
  }
}

export default App
