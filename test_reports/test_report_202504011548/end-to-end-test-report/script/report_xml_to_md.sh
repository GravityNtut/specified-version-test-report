#!/bin/bash
[ "${0##*/}" = "$0" ] || cd "${0%/*}"

cd ..

REPORT_DIR="./test_report"
BASE_URL="https://github.com/BrobridgeOrg/End-to-End-test/tree/main"

find . -type d -name "*test" -not -path "./test_report*" | while read -r dir; do
  RELATIVE_PATH="${dir#./}"
  SUB_DIR="$REPORT_DIR/$RELATIVE_PATH"

  REPORT_XML="$dir/report.xml"
  REPORT_MD="$SUB_DIR/report.md"
  mkdir -p "$SUB_DIR"

  if [ -f "$REPORT_XML" ]; then
    echo "Processing $REPORT_XML -> $REPORT_MD"

    {
      echo "# Test Report"
      echo ""
      echo "**Total tests:** $(xmlstarlet sel -t -v "/testsuites/@tests" "$REPORT_XML")"
      echo "**Failures:** $(xmlstarlet sel -t -v "/testsuites/@failures" "$REPORT_XML")"
      echo "**Errors:** $(xmlstarlet sel -t -v "/testsuites/@errors" "$REPORT_XML")"
      echo "**Time:** $(printf "%.2f" "$(xmlstarlet sel -t -v "/testsuites/@time" "$REPORT_XML")") seconds"
      echo ""

      xmlstarlet sel -t -m "/testsuites/testsuite" -v "@name" -n "$REPORT_XML" | while read -r TESTSUITE_NAME; do
        echo "## $TESTSUITE_NAME"
        echo ""

        PREV_SCENARIO=""
        FEATURE_LINE=""

        xpath_expr=$(printf 'concat(@name, "|", @status, "|", @time, "|", translate(failure/@message, "\n", "¤"))')
        xmlstarlet sel -t -m "/testsuites/testsuite[@name='$TESTSUITE_NAME']/testcase" \
          -v "$xpath_expr" -n "$REPORT_XML" \
        | sed 's/¤/<br>/g' \
        | while IFS="|" read -r TEST_NAME STATUS TEST_TIME FAILURE_MESSAGE; do
          
          SCENARIO_NAME=$(echo "$TEST_NAME" | sed 's/ [ME][0-9]*(.*\| #[0-9]*.*//')

          if [[ "$PREV_SCENARIO" != "$SCENARIO_NAME" ]]; then
            PREV_SCENARIO="$SCENARIO_NAME"
            FEATURE_FILE=""
            START_LINE=""

            for FILE in "$dir"/*.feature; do
              if [ -f "$FILE" ]; then
                ESCAPED_SCENARIO_NAME=$(echo "$SCENARIO_NAME" | sed 's/[]\/$*.^()[{}|+?]/\\&/g')
                LINE=$(grep -En "Scenario Outline: $ESCAPED_SCENARIO_NAME|Scenario: $ESCAPED_SCENARIO_NAME" "$FILE" | head -n 1 | cut -d: -f1)
                if [ -n "$LINE" ]; then
                  FEATURE_FILE="$FILE"
                  START_LINE="$LINE"
                  break
                fi
              fi
            done

            if [ -n "$FEATURE_FILE" ] && [ -n "$START_LINE" ]; then
              FEATURE_NAME=$(basename "$FEATURE_FILE")
              TRIMMED_DIR="${dir#./}"
              FEATURE_URL="$BASE_URL/$TRIMMED_DIR/$FEATURE_NAME#L${START_LINE}"
              FEATURE_LINE="[$SCENARIO_NAME]($FEATURE_URL)"
              echo ""
              echo "### $FEATURE_LINE"
              echo ""
              echo "| Test Case | Status | Time (s) | Failure Reason |"
              echo "|-----------|--------|----------|----------------|"
            fi
          fi
          
          if [[ "$TEST_NAME" == *"[waived]"* || "$TEST_NAME" == *"[Waived]"* ]]; then
            STATUS="⚠️Waived"
            TEST_NAME=${TEST_NAME//"[waived]"/}
            TEST_NAME=${TEST_NAME//"[Waived]"/}
          elif [[ "$TEST_NAME" == *"[skipped]"* ||  "$TEST_NAME" == *"[Skipped]"* ]]; then
            STATUS="⏭️Skipped"
            FAILURE_MESSAGE=""
            TEST_NAME=${TEST_NAME//"[skipped]"/}
            TEST_NAME=${TEST_NAME//"[Skipped]"/}
          elif [[ "$STATUS" == "passed" ]]; then
            STATUS="✅Passed"
          else
            STATUS="❌Failed"
          fi
          
          TEST_TIME=$(printf "%.2f" "$TEST_TIME")
          echo "| $TEST_NAME | $STATUS | $TEST_TIME | $FAILURE_MESSAGE |"
        done
        
        echo ""
      done
    } > "$REPORT_MD"

    cp "$REPORT_XML" "$SUB_DIR"
  else
    echo "Skipping $dir as no report.xml was found"
  fi

done
