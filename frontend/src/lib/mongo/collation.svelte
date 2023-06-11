<script>
  import { locales } from '$lib/mongo';

  const defaultCollation = {
    locale: 'en_US',
    strength: 3,
    caseLevel: false,
    caseFirst: 'off',
    numericOrdering: false,
    alternate: 'non-ignorable',
    maxVariable: 'punct',
    backwards: false,
    normalization: false,
  };

  export let collation = { ...defaultCollation };
</script>

<div class="settinggrid">
  <label for="collationLocale">Locale</label>
  <div class="field">
    <select id="collationLocale" bind:value={collation.locale}>
      {#each Object.entries(locales) as [ value, name ]}
        <option {value}>({value}) {name}</option>
      {/each}
    </select>
  </div>

  <label for="collationStrength">Strength</label>
  <div class="field">
    <select id="collationStrength" bind:value={collation.strength}>
      <option value={1}>(1) Collate base characters only.</option>
      <option value={2}>(2) Collate up to secondary differences, such as diacritics.</option>
      <option value={3}>(3) Collate up to tertiary differences, such as case and letter variants.</option>
      <option value={4}>(4) Limited for specific use case such as processing Japanese text.</option>
      <option value={5}>(5) Identical Level. Limited for specific use case of tie breaker.</option>
    </select>
  </div>

  <label for="collationCaseLevel">Use case level</label>
  <div class="field">
    <span class="checkbox">
      <input type="checkbox" id="collationCaseLevel" bind:checked={collation.caseLevel} />
    </span>
  </div>

  <label for="collationCaseFirst">Case first</label>
  <div class="field">
    <select id="collationCaseFirst" bind:value={collation.caseFirst}>
      <option value="off">(off) Similar to "lower" with slight differences.</option>
      <option value="upper">(upper) Uppercase sorts before lowercase.</option>
      <option value="lower">(lower) Lowercase sorts before uppercase.</option>
    </select>
  </div>

  <label for="collationNumericOrdering">Numeric ordering</label>
  <div class="field">
    <span class="checkbox">
      <input type="checkbox" id="collationNumericOrdering" bind:checked={collation.numericOrdering} />
    </span>
  </div>

  <label for="collationAlternate">Alternate</label>
  <div class="field">
    <select id="collationAlternate" bind:value={collation.alternate}>
      <option value="non-ignorable">(non-ignorable) Whitespace and punctuation are considered base characters.</option>
      <option value="shifted">(shifted) Whitespace and punctuation are considered base characters.</option>
    </select>
  </div>

  <label for="collationMaxVariable">Max Variable</label>
  <div class="field">
    <select id="collationMaxVariable" bind:value={collation.maxVariable}>
      <option value="punct">(punct) Both whitespace and punctuation are ignorable and not considered base characters.</option>
      <option value="space">(space) Whitespace is ignorable and not considered to be base characters.</option>
    </select>
  </div>

  <label for="collationBackwards">Backwards</label>
  <div class="field">
    <span class="checkbox">
      <input type="checkbox" id="collationBackwards" bind:checked={collation.backwards} />
    </span>
  </div>

  <label for="collationNormalization">Normalization</label>
  <div class="field">
    <span class="checkbox">
      <input type="checkbox" id="collationNormalization" bind:checked={collation.normalization} />
    </span>
  </div>
</div>
