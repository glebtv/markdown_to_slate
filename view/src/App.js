import React, { Component } from 'react'
import {debounce } from 'lodash'
import axios from 'axios'
import MarkdownInput from './MarkdownInput'
import SlatePreview from './SlatePreview'
import JSONTree from 'react-json-tree'
import './App.css'

class App extends Component {
  state = {
    md: "",
    slate: [],
    queue: [],
    pending: false,
    examples: [],
    example: "",
    expandAll: false,
  }

  endpoint() {
    return "http://localhost:8080"
  }


  setSlate = debounce((value) => {
    if (this.state.pending) {
      return this.setState({queue: [value]})
    }
    this.setState({pending: true})
    axios.post(this.endpoint() + "/md2slate", {
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

  setExample = event => {
    this.setState({example: event.target.value})
    if (event.target.value === "") return
    axios.get(this.endpoint() + event.target.value).then(resp => {
      this.setMD(resp.data)
    }).catch(err => {
      this.setState({error: err.message})
    })
  }

  componentDidMount() {
    axios.get(this.endpoint() + "/examples").then(resp => {
      this.setState({examples: ["", ...resp.data]})
    }).catch(err => {
      this.setState({error: err.message})
    })
  }

  expandFunc = () => {
    switch(this.state.expand) {
      case "all":
        return ()=>true
      case "none":
        return ()=>false
      default:
        return null
    }
  }

  renderExpandButtons() {
    return [
      <button key="all" onClick={()=>this.setState({expand: 'all'})}>Развернуть</button>,
      <button key="none" onClick={()=>this.setState({expand: 'none'})}>Свернуть</button>
    ]
  }

  render() {
    return (
      <div className="App">
        <MarkdownInput value={this.state.md} onChange={this.setMD} />
        <SlatePreview value={this.state.slate} />
        <div>
          <select style={{width: '100%'}} value={this.state.example} onChange={this.setExample}>
            {this.state.examples.map((ex) => (<option key={ex} value={ex}>{ex}</option>))}
          </select>
          <div>{this.renderExpandButtons()}</div>
        </div>
        <div>{this.state.error}</div>
        <JSONTree
          data={this.state.slate}
          shouldExpandNode={this.expandFunc()}
        />
      </div>
    )
  }
}

export default App
