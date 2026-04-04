// Copyright 2022 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

import ReactDomServer from "react-dom/server";
import { getHighlightedText } from "./highlightedText";

function elementToString(value: string | JSX.Element): string {
  if (typeof value == "string") {
    return value;
  }
  return ReactDomServer.renderToString(value);
}

describe("Highlighted Text", () => {
  it("text with no highlight", () => {
    const highlightedText = getHighlightedText(
      "full text",
      "no matches",
      false,
      true,
    );
    expect(highlightedText).toEqual(["full text"]);
  });

  it("text with everything highlighted", () => {
    const highlightedText = getHighlightedText(
      "everything matches",
      "everything matches",
      false,
      true,
    );

    expect(highlightedText.length).toEqual(5);
    expect(elementToString(highlightedText[1])).toEqual(
      '<span class="_text-bold" data-reactroot="">everything</span>',
    );
    expect(elementToString(highlightedText[3])).toEqual(
      '<span class="_text-bold" data-reactroot="">matches</span>',
    );
  });

  it("text with partial highlight", () => {
    const highlightedText = getHighlightedText(
      "regular text highlighted match",
      "highlighted match",
      false,
      true,
    );
    expect(highlightedText.length).toEqual(5);
    expect(highlightedText[0].toString()).toEqual("regular text ");
    expect(elementToString(highlightedText[1])).toEqual(
      '<span class="_text-bold" data-reactroot="">highlighted</span>',
    );
    expect(elementToString(highlightedText[3])).toEqual(
      '<span class="_text-bold" data-reactroot="">match</span>',
    );
  });

  it("special characters (used on regex) don't get highlighted", () => {
    const highlightedText = getHighlightedText(
      "text * ? + \\ - ^ () {} matches",
      "* ? + \\ - ^ { } ( )",
      false,
      true,
    );
    expect(highlightedText.length).toEqual(1);
    expect(highlightedText[0].toString()).toEqual(
      "text * ? + \\ - ^ () {} matches",
    );
  });
});
