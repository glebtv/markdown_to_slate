import React, {Component} from 'react'

export default class MarkdownInput extends Component {
  render() {
    return <textarea value={this.props.value} onChange={(event) => this.props.onChange(event.target.value)} />
  }
}
