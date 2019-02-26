import React, {Component} from 'react'
import './index.css'

export default class MarkdownInput extends Component {
  render() {
    return (
      <textarea
        className="MarkdownInput"
        value={this.props.value}
        onChange={(event) => this.props.onChange(event.target.value)}
      />
    )
  }
}
